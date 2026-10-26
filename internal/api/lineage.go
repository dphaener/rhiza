package api

// This file holds the lineage reads: why walks back from a node to the
// reasons it is the way it is, and log lists the change sets that touched a
// node or a link.
//
// The why walk is breadth-first and visits every node at most once, so
// cycles end. It starts from the node's live decided-by links and those of
// each of its tree ancestors (the nearest first; those decisions are
// inherited, reached via the ancestor), and the node's own live
// derived-from links, which lead to sources. From each decision reached it
// follows, in link ID order, the decision's live decided-by links out (to
// the decisions it rests on), its live supersedes links out (the decisions
// it replaced) and in (the decisions that replaced it), and its live
// derived-from links out (its sources); the questions it answers (its live
// answers links out) are listed with it. A queried node that is itself a
// decision is walked the same way, as the starting point rather than as a
// reason, and its own fields lead the answer. Retired links are skipped, and no other link type is followed. The
// far end of a link is listed whatever its kind, but only decisions are
// walked further. Then come the distinct change sets of the node's history,
// in apply order, with the sources their open names and those their own
// live derived-from links lead to.
//
// The log is the target's history: each distinct change set in apply order,
// with the operations it applied to the target, and for a landed one the
// time and actor of its land entry. A change set's carve entries that
// moved the target are listed with it, naming the other side: the target
// change set's entry re-created some of the target's operations, and the
// source's entry, found through it, dropped them. A source missing from the
// history (its operations were all dropped) is listed just before the
// change set it was carved into, with no operations, and is followed back
// the same way when it was itself carved into.
//
// Both read a bounded number of rows per node visited and the entries of
// the change sets in the target's history only.

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"slices"
	"sort"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

func init() {
	register(Operation{
		Name:     "why",
		Summary:  "Walk back through the decisions, sources and change sets that shaped a node",
		Request:  reflect.TypeFor[WhyRequest](),
		Response: reflect.TypeFor[WhyResponse](),
		Run:      runWhy,
	})
	register(Operation{
		Name:     "log",
		Summary:  "List every change set that touched a node or link, in order",
		Request:  reflect.TypeFor[LogRequest](),
		Response: reflect.TypeFor[LogResponse](),
		Run:      runLog,
	})
}

// The relations a decision is reached by.
const (
	whyDecidedBy    = "decided by"
	whySupersedes   = "supersedes"
	whySupersededBy = "superseded by"
)

// lineageAtLayout is how a land time is spelled: as the entry spells it.
const lineageAtLayout = "2006-01-02T15:04:05Z"

// WhyRequest names the node to explain.
type WhyRequest struct {
	Ref  string `json:"ref" arg:"0" desc:"The node: an ID, a unique ID prefix of 4+ characters, or a handle"`
	View string `json:"view,omitempty" desc:"as-is, to-be or what-if:<change set>; default what-if of a draft or proposed effective change set, to-be of an accepted one, else as-is"`
}

// WhyResponse is the node's reasons, sources and change sets, in walk order.
type WhyResponse struct {
	Node       WhyNode        `json:"node" desc:"The node explained"`
	Decisions  []WhyDecision  `json:"decisions" desc:"The decisions reached, in walk order"`
	Sources    []WhySource    `json:"sources" desc:"The sources reached, in walk order"`
	ChangeSets []WhyChangeSet `json:"change_sets" desc:"The change sets of the node's history, in apply order"`
	View       string         `json:"view" desc:"The view read"`
}

// WhyNode is the node explained, with every field: asked about a decision,
// why gives back its own words first.
type WhyNode struct {
	ID     id.ID                      `json:"id" desc:"The node's ID"`
	Handle string                     `json:"handle" desc:"Its handle"`
	Kind   op.Kind                    `json:"kind" desc:"Its kind"`
	Fields map[string]json.RawMessage `json:"fields" desc:"Every field, name to JSON value"`
}

// WhyRef names a node briefly.
type WhyRef struct {
	ID     id.ID  `json:"id" desc:"The node's ID"`
	Handle string `json:"handle" desc:"Its handle"`
}

// WhyDecision is one decision reached, with every field.
type WhyDecision struct {
	ID        id.ID                      `json:"id" desc:"The decision's ID"`
	Handle    string                     `json:"handle" desc:"Its handle"`
	Kind      op.Kind                    `json:"kind" desc:"Its kind: decision, unless the model links something else"`
	Fields    map[string]json.RawMessage `json:"fields" desc:"Every field, name to JSON value"`
	Via       WhyRef                     `json:"via" desc:"The node whose link led here"`
	Relation  string                     `json:"relation" enum:"decided by,supersedes,superseded by" desc:"How via relates to it: via is decided by it, via supersedes it, or via is superseded by it"`
	Inherited bool                       `json:"inherited" desc:"True when via is a tree ancestor of the node, not the node or a decision"`
	Answers   []WhyRef                   `json:"answers" desc:"The questions it answers"`
}

// WhySource is one source reached, with every field.
type WhySource struct {
	ID     id.ID                      `json:"id" desc:"The source's ID"`
	Handle string                     `json:"handle" desc:"Its handle"`
	Kind   op.Kind                    `json:"kind" desc:"Its kind: source, unless the model links something else"`
	Fields map[string]json.RawMessage `json:"fields" desc:"Every field, name to JSON value"`
	Via    WhyRef                     `json:"via" desc:"The node whose derived-from link led here"`
}

// WhyChangeSet is one change set of the node's history.
type WhyChangeSet struct {
	ID      id.ID           `json:"id" desc:"The change set's ID"`
	State   changeset.State `json:"state" desc:"Its lifecycle state"`
	Kind    changeset.Kind  `json:"kind" desc:"intention or record"`
	Intent  string          `json:"intent" desc:"Its intent"`
	Sources []id.ID         `json:"sources" desc:"The sources its open names, then those its own derived-from links lead to"`
}

// whyWalk is the state of one why walk.
type whyWalk struct {
	r    *reader
	seen map[id.ID]bool
	// queue holds the decisions still to walk from.
	queue []gmodel.Node
	resp  *WhyResponse
	// answers indexes resp.Decisions by ID.
	answers map[id.ID]int
}

func runWhy(ctx *Ctx, req any) (any, error) {
	q := req.(WhyRequest)
	r, err := openReader(ctx, "why", q.View, defaultView)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	n, err := r.Resolve(q.Ref)
	if err != nil {
		return nil, err
	}
	resp := &WhyResponse{
		Node:       WhyNode{ID: n.ID, Handle: n.Handle, Kind: n.Kind, Fields: lineageFields(n.Fields)},
		Decisions:  []WhyDecision{},
		Sources:    []WhySource{},
		ChangeSets: []WhyChangeSet{},
		View:       r.View(),
	}
	w := &whyWalk{r: r, seen: map[id.ID]bool{n.ID: true}, resp: resp, answers: map[id.ID]int{}}
	if err := w.start(n); err != nil {
		return nil, err
	}
	for len(w.queue) > 0 {
		d := w.queue[0]
		w.queue = w.queue[1:]
		if err := w.decision(d); err != nil {
			return nil, err
		}
	}
	css, err := lineageChangeSets(r, n.History)
	if err != nil {
		return nil, err
	}
	for _, c := range css {
		srcs, err := whyChangeSetSources(r, c)
		if err != nil {
			return nil, err
		}
		resp.ChangeSets = append(resp.ChangeSets, WhyChangeSet{ID: c.info.ID, State: c.info.State,
			Kind: c.info.Kind, Intent: c.info.Intent, Sources: srcs})
	}
	return *resp, nil
}

// start walks the node's own links and its ancestors' decided-by links.
func (w *whyWalk) start(n gmodel.Node) error {
	ancestors, err := w.r.Ancestors(n.ID)
	if err != nil {
		return err
	}
	if n.Kind == op.KindDecision {
		// The starting point: walked as a decision, not listed as one.
		if err := w.decision(n); err != nil {
			return err
		}
	} else {
		ls, err := w.r.Links(n.ID, index.From)
		if err != nil {
			return err
		}
		for _, l := range ls {
			if l.Retired {
				continue
			}
			switch l.Type {
			case op.LinkDecidedBy:
				if err := w.addDecision(l.To, n, whyDecidedBy, false); err != nil {
					return err
				}
			case op.LinkDerivedFrom:
				if err := w.addSource(l.To, n); err != nil {
					return err
				}
			}
		}
	}
	slices.Reverse(ancestors)
	for _, a := range ancestors {
		ls, err := w.r.Links(a.ID, index.From)
		if err != nil {
			return err
		}
		for _, l := range ls {
			if !l.Retired && l.Type == op.LinkDecidedBy {
				if err := w.addDecision(l.To, a, whyDecidedBy, true); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// decision walks the links of decision d: decided-by, supersedes both ways,
// derived-from and answers, in link ID order. The answers go on the listed
// decision of d, if it is listed.
func (w *whyWalk) decision(d gmodel.Node) error {
	ls, err := w.r.Links(d.ID, index.Both)
	if err != nil {
		return err
	}
	for _, l := range ls {
		if l.Retired {
			continue
		}
		out := l.From == d.ID
		var err error
		switch {
		case out && l.Type == op.LinkDecidedBy:
			err = w.addDecision(l.To, d, whyDecidedBy, false)
		case out && l.Type == op.LinkSupersedes:
			err = w.addDecision(l.To, d, whySupersedes, false)
		case !out && l.Type == op.LinkSupersedes:
			err = w.addDecision(l.From, d, whySupersededBy, false)
		case out && l.Type == op.LinkDerivedFrom:
			err = w.addSource(l.To, d)
		case out && l.Type == op.LinkAnswers:
			err = w.addAnswer(d.ID, l.To)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// addDecision lists node i, reached from via by relation, unless it was
// visited; a decision is queued to walk from.
func (w *whyWalk) addDecision(i id.ID, via gmodel.Node, relation string, inherited bool) error {
	if w.seen[i] {
		return nil
	}
	n, ok, err := w.r.Node(i)
	if err != nil || !ok {
		return err
	}
	w.seen[i] = true
	w.answers[i] = len(w.resp.Decisions)
	w.resp.Decisions = append(w.resp.Decisions, WhyDecision{ID: n.ID, Handle: n.Handle, Kind: n.Kind,
		Fields: lineageFields(n.Fields), Via: WhyRef{ID: via.ID, Handle: via.Handle},
		Relation: relation, Inherited: inherited, Answers: []WhyRef{}})
	if n.Kind == op.KindDecision {
		w.queue = append(w.queue, n)
	}
	return nil
}

// addSource lists node i, reached from via by derived-from, unless it was
// visited.
func (w *whyWalk) addSource(i id.ID, via gmodel.Node) error {
	if w.seen[i] {
		return nil
	}
	n, ok, err := w.r.Node(i)
	if err != nil || !ok {
		return err
	}
	w.seen[i] = true
	w.resp.Sources = append(w.resp.Sources, WhySource{ID: n.ID, Handle: n.Handle, Kind: n.Kind,
		Fields: lineageFields(n.Fields), Via: WhyRef{ID: via.ID, Handle: via.Handle}})
	return nil
}

// addAnswer adds node q to what decision d answers, when d is listed.
func (w *whyWalk) addAnswer(d, q id.ID) error {
	k, listed := w.answers[d]
	if !listed {
		return nil
	}
	n, ok, err := w.r.Node(q)
	if err != nil || !ok {
		return err
	}
	w.resp.Decisions[k].Answers = append(w.resp.Decisions[k].Answers, WhyRef{ID: n.ID, Handle: n.Handle})
	return nil
}

// lineageFields is fs, or an empty map for nil, so fields encode as {}.
func lineageFields(fs map[string]json.RawMessage) map[string]json.RawMessage {
	if fs == nil {
		return map[string]json.RawMessage{}
	}
	return fs
}

// lineageCS is one change set of a history, with its entries and the
// operations of the history it applied.
type lineageCS struct {
	info    index.ChangeSetInfo
	entries []entry.Entry
	ops     []id.ID
}

// lineageChangeSets groups history by change set, in order of first
// appearance, reading each change set's entries once.
func lineageChangeSets(r *reader, history []gmodel.Ref) ([]lineageCS, error) {
	var out []lineageCS
	at := map[id.ID]int{}
	for _, h := range history {
		if k, ok := at[h.ChangeSet]; ok {
			out[k].ops = append(out[k].ops, h.Op)
			continue
		}
		at[h.ChangeSet] = len(out)
		out = append(out, lineageCS{info: index.ChangeSetInfo{ID: h.ChangeSet}, ops: []id.ID{h.Op}})
	}
	if len(out) == 0 {
		return nil, nil
	}
	infos, err := r.ChangeSets()
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		if k, ok := at[info.ID]; ok {
			out[k].info = info
		}
	}
	for k := range out {
		es, err := r.Entries(out[k].info.ID)
		if err != nil {
			return nil, err
		}
		out[k].entries = es
	}
	return out, nil
}

// lineageEntry is the entry of c with ID e.
func lineageEntry(c lineageCS, e id.ID) (entry.Entry, bool) {
	if e.IsZero() {
		return entry.Entry{}, false
	}
	k := sort.Search(len(c.entries), func(k int) bool { return id.Compare(c.entries[k].ID, e) >= 0 })
	if k < len(c.entries) && c.entries[k].ID == e {
		return c.entries[k], true
	}
	return entry.Entry{}, false
}

// whyChangeSetSources is the sources c's open names, then those its own
// live derived-from links lead to.
func whyChangeSetSources(r *reader, c lineageCS) ([]id.ID, error) {
	out := []id.ID{}
	if e, ok := lineageEntry(c, c.info.OpenEntry); ok {
		if b, err := changeset.DecodeBody(e); err == nil {
			if open, ok := b.(changeset.Open); ok {
				out = append(out, open.Sources...)
			}
		}
	}
	ls, err := r.Links(c.info.ID, index.From)
	if err != nil {
		return nil, err
	}
	for _, l := range ls {
		if !l.Retired && l.Type == op.LinkDerivedFrom && !slices.Contains(out, l.To) {
			out = append(out, l.To)
		}
	}
	return out, nil
}

// Human prints the node's own fields, then the decisions and sources with
// every field in full (they are the owner's words), then the change sets.
func (r WhyResponse) Human(w io.Writer) {
	fmt.Fprintf(w, "# Why %s\n\n", r.Node.Handle)
	if len(r.Node.Fields) > 0 {
		fmt.Fprintf(w, "- %s %s\n", r.Node.Kind, r.Node.ID)
		lineageHumanFields(w, r.Node.Fields)
	}
	fmt.Fprintf(w, "## Decisions\n\n")
	if len(r.Decisions) == 0 {
		fmt.Fprintf(w, "None.\n\n")
	}
	for _, d := range r.Decisions {
		fmt.Fprintf(w, "### %s\n\n", d.Handle)
		line := fmt.Sprintf("- %s %s %s", d.Via.Handle, d.Relation, d.Handle)
		if d.Inherited {
			line += " (inherited)"
		}
		fmt.Fprintf(w, "%s\n- ID: %s\n", line, d.ID)
		if d.Kind != op.KindDecision {
			fmt.Fprintf(w, "- Kind: %s\n", d.Kind)
		}
		for _, a := range d.Answers {
			fmt.Fprintf(w, "- Answers: %s (%s)\n", a.Handle, shortID(a.ID))
		}
		lineageHumanFields(w, d.Fields)
	}
	fmt.Fprintf(w, "## Sources\n\n")
	if len(r.Sources) == 0 {
		fmt.Fprintf(w, "None.\n\n")
	}
	for _, s := range r.Sources {
		fmt.Fprintf(w, "### %s\n\n- %s derived from %s\n- ID: %s\n", s.Handle, s.Via.Handle, s.Handle, s.ID)
		if s.Kind != op.KindSource {
			fmt.Fprintf(w, "- Kind: %s\n", s.Kind)
		}
		lineageHumanFields(w, s.Fields)
	}
	fmt.Fprintf(w, "## Change sets\n\n")
	if len(r.ChangeSets) == 0 {
		fmt.Fprintln(w, "None.")
	}
	for _, c := range r.ChangeSets {
		line := fmt.Sprintf("- %s %s %s: %s", c.ID, c.State, c.Kind, c.Intent)
		if len(c.Sources) > 0 {
			line += " (sources: " + joinIDs(c.Sources) + ")"
		}
		fmt.Fprintln(w, line)
	}
}

// lineageHumanFields prints every field in name order, its text in full,
// then a blank line.
func lineageHumanFields(w io.Writer, fs map[string]json.RawMessage) {
	names := make([]string, 0, len(fs))
	for name := range fs {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		fmt.Fprintf(w, "- %s: %s\n", name, fieldText(fs[name]))
	}
	fmt.Fprintln(w)
}

// LogRequest names the node or link whose history to list.
type LogRequest struct {
	Ref  string `json:"ref" arg:"0" desc:"The node or link: an ID, a unique ID prefix of 4+ characters, or a node's handle"`
	View string `json:"view,omitempty" desc:"as-is, to-be or what-if:<change set>; default what-if of a draft or proposed effective change set, to-be of an accepted one, else as-is"`
}

// LogResponse is the target's history, one change set at a time.
type LogResponse struct {
	Target     LogTarget      `json:"target" desc:"The node or link"`
	ChangeSets []LogChangeSet `json:"change_sets" desc:"Every change set that touched it, in apply order"`
	View       string         `json:"view" desc:"The view read"`
}

// LogTarget names the node or link logged.
type LogTarget struct {
	ID     id.ID  `json:"id" desc:"Its ID"`
	Kind   string `json:"kind" enum:"node,link" desc:"node or link"`
	Handle string `json:"handle,omitempty" desc:"The node's handle; absent for a link"`
}

// LogChangeSet is one change set of the target's history.
type LogChangeSet struct {
	ID     id.ID           `json:"id" desc:"The change set's ID"`
	State  changeset.State `json:"state" desc:"Its lifecycle state"`
	Kind   changeset.Kind  `json:"kind" desc:"intention or record"`
	Intent string          `json:"intent" desc:"Its intent"`
	Ops    []LogOp         `json:"ops" desc:"The operations it applied to the target, in apply order"`
	At     string          `json:"at,omitempty" desc:"When it landed (RFC 3339, UTC); absent unless landed"`
	Actor  id.ID           `json:"actor,omitzero" desc:"Who landed it; absent unless landed"`
	Carves []LogCarve      `json:"carves,omitempty" desc:"Its carve entries that moved the target in or out, in the order found; absent when none"`
	Events []LogEvent      `json:"events,omitempty" desc:"Its applied land and reopen entries, in causal order; absent unless it was ever reopened"`
}

// LogEvent is one applied land or reopen of a change set.
type LogEvent struct {
	Type  entry.Type `json:"type" desc:"land or reopen"`
	Entry id.ID      `json:"entry" desc:"The land or reopen entry"`
	At    string     `json:"at" desc:"When it was written (RFC 3339, UTC)"`
	Actor id.ID      `json:"actor" desc:"Who landed or reopened it"`
}

// LogCarve is one carve entry of a change set that moved the target. On the
// target change set's entry it names the source, on the source's the target.
type LogCarve struct {
	Entry      id.ID `json:"entry" desc:"The carve entry, in this change set"`
	CarvedFrom id.ID `json:"carved_from,omitzero" desc:"The change set the target was carved from; set on the target change set's entry"`
	CarvedInto id.ID `json:"carved_into,omitzero" desc:"The change set the target was carved into; set on the source change set's entry"`
	CarveEntry id.ID `json:"carve_entry" desc:"The other change set's entry of the same carve"`
}

// LogOp is one operation applied to the target.
type LogOp struct {
	ID   id.ID   `json:"id" desc:"The operation's ID"`
	Type op.Type `json:"type" desc:"Its type, such as node.create"`
}

func runLog(ctx *Ctx, req any) (any, error) {
	q := req.(LogRequest)
	r, err := openReader(ctx, "log", q.View, defaultView)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var target LogTarget
	var history []gmodel.Ref
	n, err := r.Resolve(q.Ref)
	switch {
	case err == nil:
		target, history = LogTarget{ID: n.ID, Kind: "node", Handle: n.Handle}, n.History
	case AsError(err).Code == CodeNotFound:
		l, lerr := r.ResolveLink(q.Ref)
		if lerr != nil {
			if AsError(lerr).Code == CodeNotFound {
				return nil, NotFoundf("no node or link matches %q", q.Ref).WithNext(lsNext)
			}
			return nil, lerr
		}
		target, history = LogTarget{ID: l.ID, Kind: "link"}, l.History
	default:
		return nil, err
	}
	css, err := lineageChangeSets(r, history)
	if err != nil {
		return nil, err
	}
	css, carves, err := logCarves(r, css)
	if err != nil {
		return nil, err
	}
	resp := LogResponse{Target: target, ChangeSets: []LogChangeSet{}, View: r.View()}
	for _, c := range css {
		lc := LogChangeSet{ID: c.info.ID, State: c.info.State, Kind: c.info.Kind, Intent: c.info.Intent, Ops: []LogOp{},
			Carves: carves[c.info.ID]}
		folded, _ := changeset.Fold(c.info.ID, c.entries)
		types := logOpTypes(folded)
		for _, o := range c.ops {
			lc.Ops = append(lc.Ops, LogOp{ID: o, Type: types[o]})
		}
		lc.Events = logEvents(folded)
		if c.info.State == changeset.StateLanded {
			if e, ok := lineageEntry(c, c.info.TerminalEntry); ok {
				lc.At, lc.Actor = e.At.UTC().Format(lineageAtLayout), e.Actor
			}
		}
		resp.ChangeSets = append(resp.ChangeSets, lc)
	}
	return resp, nil
}

// logCarveWalk is the state of finding the carves of one log.
type logCarveWalk struct {
	r *reader
	// listed holds the change sets already in the log.
	listed map[id.ID]bool
	// infos is every change set's info, read when a source must be listed.
	infos map[id.ID]index.ChangeSetInfo
	// added holds the sources to list before the change set being walked.
	added  []lineageCS
	carves map[id.ID][]LogCarve
	seen   map[id.ID]bool
}

// logCarves returns css with the sources the target was carved out of
// added, each before the change set it was carved into, and each change
// set's carve entries that moved the target. Only entries of format 2 or
// later can carry a carve, so a log without carves decodes nothing more.
func logCarves(r *reader, css []lineageCS) ([]lineageCS, map[id.ID][]LogCarve, error) {
	l := &logCarveWalk{r: r, listed: map[id.ID]bool{}, carves: map[id.ID][]LogCarve{}, seen: map[id.ID]bool{}}
	for _, c := range css {
		l.listed[c.info.ID] = true
	}
	var out []lineageCS
	for _, c := range css {
		l.added = nil
		if err := l.walk(c, c.ops); err != nil {
			return nil, nil, err
		}
		out = append(append(out, l.added...), c)
	}
	return out, l.carves, nil
}

// walk records the carve entries of c that re-created any of ops, and the
// source entries they name, then walks each source with the ops its entry
// dropped. A source not yet listed is added after its own sources.
func (l *logCarveWalk) walk(c lineageCS, ops []id.ID) error {
	for _, e := range c.entries {
		o, ok := logCarveBody(e)
		if !ok || o.Carve.From.IsZero() || l.seen[e.ID] ||
			!slices.ContainsFunc(o.Ops, func(x op.Op) bool { return slices.Contains(ops, x.ID) }) {
			continue
		}
		l.seen[e.ID] = true
		l.carves[c.info.ID] = append(l.carves[c.info.ID], LogCarve{Entry: e.ID, CarvedFrom: o.Carve.From, CarveEntry: o.Carve.Entry})
		src, err := l.source(o.Carve.From)
		if err != nil {
			return err
		}
		se, ok := lineageEntry(src, o.Carve.Entry)
		if !ok || l.seen[se.ID] {
			continue
		}
		so, ok := logCarveBody(se)
		if !ok || so.Carve.To != c.info.ID {
			continue
		}
		l.seen[se.ID] = true
		l.carves[src.info.ID] = append(l.carves[src.info.ID], LogCarve{Entry: se.ID, CarvedInto: c.info.ID, CarveEntry: e.ID})
		if err := l.walk(src, so.Drop); err != nil {
			return err
		}
		if !l.listed[src.info.ID] {
			l.listed[src.info.ID] = true
			l.added = append(l.added, src)
		}
	}
	return nil
}

// source is change set cs with its entries and no operations.
func (l *logCarveWalk) source(cs id.ID) (lineageCS, error) {
	if l.infos == nil {
		infos, err := l.r.ChangeSets()
		if err != nil {
			return lineageCS{}, err
		}
		l.infos = make(map[id.ID]index.ChangeSetInfo, len(infos))
		for _, info := range infos {
			l.infos[info.ID] = info
		}
	}
	info, ok := l.infos[cs]
	if !ok {
		info = index.ChangeSetInfo{ID: cs}
	}
	es, err := l.r.Entries(cs)
	if err != nil {
		return lineageCS{}, err
	}
	return lineageCS{info: info, entries: es}, nil
}

// logCarveBody is e's ops body when it carries a carve.
func logCarveBody(e entry.Entry) (changeset.Ops, bool) {
	if e.V < 2 || e.Type != entry.TypeOps {
		return changeset.Ops{}, false
	}
	b, err := changeset.DecodeBody(e)
	if err != nil {
		return changeset.Ops{}, false
	}
	o, ok := b.(changeset.Ops)
	return o, ok && o.Carve != nil
}

// logEvents is folded's applied land and reopen entries, in causal order,
// or nil when it was never reopened: a land alone is At and Actor.
func logEvents(folded changeset.ChangeSet) []LogEvent {
	var out []LogEvent
	reopened := false
	for _, m := range folded.Lifecycle {
		if m.Type != entry.TypeLand && m.Type != entry.TypeReopen {
			continue
		}
		reopened = reopened || m.Type == entry.TypeReopen
		out = append(out, LogEvent{Type: m.Type, Entry: m.Entry, At: m.At.UTC().Format(lineageAtLayout), Actor: m.Actor})
	}
	if !reopened {
		return nil
	}
	return out
}

// logOpTypes maps each operation of the folded change set to its type.
func logOpTypes(folded changeset.ChangeSet) map[id.ID]op.Type {
	out := make(map[id.ID]op.Type, len(folded.Ops))
	for _, o := range folded.Ops {
		out[o.ID] = o.Type
	}
	return out
}

// Human prints one block per change set.
func (r LogResponse) Human(w io.Writer) {
	name := r.Target.Handle
	if r.Target.Kind == "link" {
		name = "link " + r.Target.ID.String()
	}
	fmt.Fprintf(w, "# Log %s\n", name)
	if len(r.ChangeSets) == 0 {
		fmt.Fprintf(w, "\nNo change sets.\n")
	}
	for _, c := range r.ChangeSets {
		fmt.Fprintf(w, "\n## %s %s %s: %s\n\n", c.ID, c.State, c.Kind, c.Intent)
		for _, e := range c.Events {
			verb := "Landed"
			if e.Type == entry.TypeReopen {
				verb = "Reopened"
			}
			fmt.Fprintf(w, "- %s %s by %s\n", verb, e.At, e.Actor)
		}
		if c.At != "" && len(c.Events) == 0 {
			fmt.Fprintf(w, "- Landed %s by %s\n", c.At, c.Actor)
		}
		for _, k := range c.Carves {
			from, into := k.CarvedFrom, c.ID
			if from.IsZero() {
				from, into = c.ID, k.CarvedInto
			}
			fmt.Fprintf(w, "- Carved from %s into %s\n", shortID(from), shortID(into))
		}
		for _, o := range c.Ops {
			fmt.Fprintf(w, "- %s %s\n", o.Type, o.ID)
		}
	}
}
