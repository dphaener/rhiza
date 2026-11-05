package index

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// metaAsIs is the meta key holding what the as-is view knew when it was
// built beyond its rows (see asIsBlob).
const metaAsIs = "asis_state"

// asIsBlob is the as-is state kept outside the node and link rows: the
// count of operations, the conflicts later writes may still repair, and the
// IDs landed create ops name that the as-is model lacks because the create
// was skipped.
type asIsBlob struct {
	Seq   int        `json:"seq"`
	Held  []heldBlob `json:"held,omitempty"`
	Extra []string   `json:"extra,omitempty"`
}

type heldBlob struct {
	Target, ChangeSet, Op string
	Left                  []string
}

// writeAsIsState stores the as-is state st of model m over sets.
func writeAsIsState(tx *sql.Tx, sets []changeset.ChangeSet, m model.Model, st model.State) error {
	b := asIsBlob{Seq: st.Seq}
	for _, h := range st.Held {
		b.Held = append(b.Held, heldBlob{idText(h.Target), idText(h.ChangeSet), idText(h.Op), h.Left})
	}
	for _, s := range sets {
		if s.State != changeset.StateLanded {
			continue
		}
		for _, o := range s.Ops {
			if o.Type != op.TypeNodeCreate && o.Type != op.TypeLinkCreate {
				continue
			}
			t := o.Target()
			if _, ok := m.Node(t); ok {
				continue
			}
			if _, ok := m.Link(t); ok {
				continue
			}
			b.Extra = append(b.Extra, idText(t))
		}
	}
	slices.Sort(b.Extra)
	b.Extra = slices.Compact(b.Extra)
	data, _ := json.Marshal(b)
	_, err := tx.Exec(`INSERT OR REPLACE INTO meta (key, value) VALUES (?, ?)`, metaAsIs, string(data))
	return err
}

// encodeAspects spells the aspects of a history as JSON; nil (SQL NULL)
// when there are none.
func encodeAspects(a [][]string) any {
	if len(a) == 0 {
		return nil
	}
	b, _ := json.Marshal(a)
	return string(b)
}

func decodeAspects(s sql.NullString) ([][]string, error) {
	if !s.Valid {
		return nil, nil
	}
	var a [][]string
	if err := json.Unmarshal([]byte(s.String), &a); err != nil {
		return nil, fmt.Errorf("stored aspects: %w", err)
	}
	return a, nil
}

// storedErr is a stored problem's error: its text, wrapping its sentinel.
type storedErr struct {
	msg string
	is  error
}

func (e storedErr) Error() string { return e.msg }
func (e storedErr) Unwrap() error { return e.is }

// storedError is the error of a stored problem: the sentinel is wrapped
// as the model wraps it ("sentinel: detail"), so it is the error the model
// made; anything else keeps its text and sentinel.
func storedError(msg string, sentinel error) error {
	if rest, ok := strings.CutPrefix(msg, sentinel.Error()+": "); ok {
		return fmt.Errorf("%w: %s", sentinel, rest)
	}
	return storedErr{msg, sentinel}
}

// loadSeed reads the stored as-is view as a model.Seed.
func loadSeed(q querier) (model.Seed, error) {
	seed, err := loadSeedBase(q)
	if err != nil {
		return seed, err
	}
	return loadSeedRows(q, seed)
}

// loadSeedRows adds every stored as-is node and link to seed.
func loadSeedRows(q querier, seed model.Seed) (model.Seed, error) {
	nrows, err := q.Query(nodeSelect[:len(nodeSelect)-len(" FROM nodes")]+`, aspects, placed FROM nodes WHERE view = ?`, viewAsIs)
	if err != nil {
		return seed, err
	}
	for nrows.Next() {
		var asp, placed sql.NullString
		n, err := scanNode(nrows, &asp, &placed)
		if err != nil {
			nrows.Close()
			return seed, err
		}
		sn := model.SeedNode{Node: n}
		if sn.Aspects, err = decodeAspects(asp); err != nil {
			nrows.Close()
			return seed, err
		}
		if parts := strings.Fields(placed.String); len(parts) == 3 {
			sn.Placed.Seq, _ = strconv.Atoi(parts[0])
			sn.Placed.Ref = model.Ref{ChangeSet: parseID(parts[1]), Op: parseID(parts[2])}
		} else if len(parts) == 2 {
			sn.Placed.Seq, _ = strconv.Atoi(parts[0])
		}
		seed.Nodes = append(seed.Nodes, sn)
	}
	nrows.Close()
	if err := nrows.Err(); err != nil {
		return seed, err
	}

	lrows, err := q.Query(linkSelect[:len(linkSelect)-len(" FROM links")]+`, aspects FROM links WHERE view = ?`, viewAsIs)
	if err != nil {
		return seed, err
	}
	for lrows.Next() {
		var asp sql.NullString
		l, err := scanLink(lrows, &asp)
		if err != nil {
			lrows.Close()
			return seed, err
		}
		sl := model.SeedLink{Link: l}
		if sl.Aspects, err = decodeAspects(asp); err != nil {
			lrows.Close()
			return seed, err
		}
		seed.Links = append(seed.Links, sl)
	}
	lrows.Close()
	if err := lrows.Err(); err != nil {
		return seed, err
	}
	return seed, nil
}

// loadSeedBase reads the stored as-is view as a model.Seed without its
// nodes and links.
func loadSeedBase(q querier) (model.Seed, error) {
	var seed model.Seed
	rows, err := q.Query(`SELECT id FROM change_sets ORDER BY id`)
	if err != nil {
		return seed, err
	}
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			rows.Close()
			return seed, err
		}
		seed.Known = append(seed.Known, parseID(s))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return seed, err
	}

	prows, err := q.Query(`SELECT change_set, op, sentinel, message FROM problems WHERE source = ? ORDER BY seq`, viewAsIs)
	if err != nil {
		return seed, err
	}
	for prows.Next() {
		var cs, o, sent, msg string
		if err := prows.Scan(&cs, &o, &sent, &msg); err != nil {
			prows.Close()
			return seed, err
		}
		var is error
		for _, s := range modelSentinels {
			if s.name == sent {
				is = s.err
			}
		}
		if is == nil || model.IsDuplicateHandle(is) {
			continue
		}
		seed.Problems = append(seed.Problems, model.Problem{ChangeSet: parseID(cs), Op: parseID(o), Err: storedError(msg, is)})
	}
	prows.Close()
	if err := prows.Err(); err != nil {
		return seed, err
	}

	var blob asIsBlob
	if v, err := metaOf(q, metaAsIs); err != nil {
		return seed, err
	} else if v != "" {
		if err := json.Unmarshal([]byte(v), &blob); err != nil {
			return seed, fmt.Errorf("stored as-is state: %w", err)
		}
	}
	seed.State.Seq = blob.Seq
	for _, h := range blob.Held {
		seed.State.Held = append(seed.State.Held, model.Held{Target: parseID(h.Target),
			ChangeSet: parseID(h.ChangeSet), Op: parseID(h.Op), Left: h.Left})
	}
	return seed, nil
}

func inFlightState(s changeset.State) bool {
	return s == changeset.StateDraft || s == changeset.StateProposed || s == changeset.StateAccepted
}

// WhatIf returns the what-if view of change set cs, as model.WhatIf
// computes it over the change sets folded from every stored entry: as-is
// with cs on top. It decodes only the entries of cs, of the in-flight
// change sets it depends on and of the landed ones it depends on, and goes
// on from the stored as-is rows. Problems are model.WhatIf's (ErrMissing
// when cs is not in the log, ErrState when it is not in flight).
func (x *Index) WhatIf(cs id.ID) (model.Model, []model.Problem, error) {
	tx, err := x.db.Begin()
	if err != nil {
		return model.Model{}, nil, err
	}
	defer tx.Rollback()

	sets, err := whatIfSets(tx, cs, nil)
	if err != nil {
		return model.Model{}, nil, err
	}
	seed, err := loadSeed(tx)
	if err != nil {
		return model.Model{}, nil, err
	}
	m, ps := model.WhatIfFrom(seed, sets, cs)
	return m, ps, nil
}

// whatIfSets returns the change sets a what-if of cs reads: every change
// set's ID and state, the whole of cs, of the in-flight change sets it
// depends on and of the landed ones it depends on, and the landed
// depends-on links, as model.WhatIfFrom wants them.
func whatIfSets(tx *sql.Tx, cs id.ID, over *changeset.ChangeSet) ([]changeset.ChangeSet, error) {
	var stubs []changeset.ChangeSet
	rows, err := tx.Query(`SELECT id, state FROM change_sets ORDER BY id`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var s, state string
		if err := rows.Scan(&s, &state); err != nil {
			rows.Close()
			return nil, err
		}
		stubs = append(stubs, changeset.ChangeSet{ID: parseID(s), State: changeset.State(state)})
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var landed []op.Op
	lrows, err := tx.Query(`SELECT retire, link, from_id, to_id FROM landed_links ORDER BY change_set, n`)
	if err != nil {
		return nil, err
	}
	for lrows.Next() {
		var retire int
		var link, from, to string
		if err := lrows.Scan(&retire, &link, &from, &to); err != nil {
			lrows.Close()
			return nil, err
		}
		if retire != 0 {
			landed = append(landed, op.Op{Type: op.TypeLinkRetire, Link: parseID(link)})
		} else {
			landed = append(landed, op.Op{Type: op.TypeLinkCreate, LinkType: op.LinkDependsOn,
				Link: parseID(link), From: parseID(from), To: parseID(to)})
		}
	}
	lrows.Close()
	if err := lrows.Err(); err != nil {
		return nil, err
	}

	full := map[id.ID]changeset.ChangeSet{}
	var extraCreates []op.Op
	if over != nil {
		// cs is the given fold, whatever the index holds of it.
		o := *over
		o.ID = cs
		full[cs] = o
		i := slices.IndexFunc(stubs, func(s changeset.ChangeSet) bool { return s.ID == cs })
		if i < 0 {
			stubs = append(stubs, changeset.ChangeSet{ID: cs})
			slices.SortFunc(stubs, func(a, b changeset.ChangeSet) int { return id.Compare(a.ID, b.ID) })
			i = slices.IndexFunc(stubs, func(s changeset.ChangeSet) bool { return s.ID == cs })
		}
		stubs[i].State = o.State
		var err error
		if extraCreates, err = createdByLanded(tx, o.Ops); err != nil {
			return nil, err
		}
	}
	build := func() []changeset.ChangeSet {
		out := make([]changeset.ChangeSet, 0, len(stubs)+1)
		for _, s := range stubs {
			if f, ok := full[s.ID]; ok {
				s = f
			}
			out = append(out, s)
		}
		ops := slices.Concat(landed, extraCreates)
		return append(out, changeset.ChangeSet{State: changeset.StateLanded, Ops: ops})
	}
	state := func(c id.ID) (changeset.State, bool) {
		for _, s := range stubs {
			if s.ID == c {
				return s.State, true
			}
		}
		return "", false
	}

	for {
		if st, ok := state(cs); !ok || !inFlightState(st) {
			break
		}
		a := changeset.Analyze(build())
		want := slices.Concat([]id.ID{cs}, a.Closure(cs), a.Dependencies(cs))
		var need []id.ID
		for _, w := range want {
			if _, ok := full[w]; !ok && !slices.Contains(need, w) {
				need = append(need, w)
			}
		}
		if len(need) == 0 {
			break
		}
		for _, w := range need {
			c, err := foldOne(tx, w)
			if err != nil {
				return nil, err
			}
			full[w] = c
		}
		if c, ok := full[cs]; ok && extraCreates == nil {
			if extraCreates, err = createdByLanded(tx, c.Ops); err != nil {
				return nil, err
			}
		}
	}
	return build(), nil
}

// foldOne folds the stored entries of change set cs.
func foldOne(q querier, cs id.ID) (changeset.ChangeSet, error) {
	rows, err := q.Query(`SELECT raw FROM entries WHERE change_set = ? ORDER BY id`, idText(cs))
	if err != nil {
		return changeset.ChangeSet{}, err
	}
	defer rows.Close()
	var es []entry.Entry
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return changeset.ChangeSet{}, err
		}
		e, err := entry.DecodeCanonical(raw)
		if err != nil {
			return changeset.ChangeSet{}, fmt.Errorf("decode stored entry: %w", err)
		}
		es = append(es, e)
	}
	if err := rows.Err(); err != nil {
		return changeset.ChangeSet{}, err
	}
	c, _ := changeset.FoldCanonical(cs, es)
	return c, nil
}

// createdByLanded returns create ops for the IDs ops name that a landed
// change set creates: the as-is model has them, or a landed create of
// them was skipped.
func createdByLanded(q querier, ops []op.Op) ([]op.Op, error) {
	var blob asIsBlob
	if v, err := metaOf(q, metaAsIs); err != nil {
		return nil, err
	} else if v != "" {
		if err := json.Unmarshal([]byte(v), &blob); err != nil {
			return nil, err
		}
	}
	var out []op.Op
	seen := map[id.ID]bool{}
	for _, o := range ops {
		for _, n := range changeset.Names(o) {
			if seen[n] {
				continue
			}
			seen[n] = true
			var one int
			err := q.QueryRow(`SELECT 1 WHERE EXISTS (SELECT 1 FROM nodes WHERE view = ? AND id = ?)
				OR EXISTS (SELECT 1 FROM links WHERE view = ? AND id = ?)`, viewAsIs, idText(n), viewAsIs, idText(n)).Scan(&one)
			if err == nil || slices.Contains(blob.Extra, idText(n)) {
				out = append(out, op.Op{Type: op.TypeNodeCreate, Node: n})
			} else if err != sql.ErrNoRows {
				return nil, err
			}
		}
	}
	if out == nil {
		out = []op.Op{}
	}
	return out, nil
}
