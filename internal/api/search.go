package api

// This file holds search: the live nodes of a view whose handle or field
// text matches, ranked, a page at a time. As-is and to-be read the index's
// FTS table for candidates, what-if scans the computed model in memory; both
// then go through the same matcher, so a view gives the same hits whichever
// backend answers. The cursor is opaque: the base64url (unpadded) form of
// "<score>:<ID>" of the last hit of the page before.

import (
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

func init() {
	register(Operation{
		Name:     "search",
		Summary:  "Find live nodes by the text of their handles and fields, ranked",
		Request:  reflect.TypeFor[SearchRequest](),
		Response: reflect.TypeFor[SearchResponse](),
		Run:      runSearch,
	})
}

// The page sizes of search.
const (
	searchDefaultLimit = 20
	searchMaxLimit     = 200
)

// searchDescribeNext is the next command for a bad search request.
const searchDescribeNext = "rhiza describe search"

// The score of one matched word, by where it is found.
const (
	scoreHandleWord   = 100
	scoreHandlePrefix = 60
	scoreFieldWord    = 20
	scoreFieldPrefix  = 10
)

// SearchRequest is the text to find and how to page the hits.
type SearchRequest struct {
	Text   string   `json:"text" arg:"0" desc:"The words to find; every word must start a word of the handle or of a field"`
	Kinds  []string `json:"kinds,omitempty" flag:"kind" enum:"intent,capability,constraint,decision,question,evidence,source,actor,anchor,note" desc:"Keep the nodes of this kind; repeat for several; default every kind"`
	Limit  int      `json:"limit,omitempty" desc:"The most hits to return, 1 to 200; default 20"`
	Cursor string   `json:"cursor,omitempty" desc:"Where to go on: the next_cursor of the page before"`
	View   string   `json:"view,omitempty" desc:"as-is, to-be or what-if:<change set>; default follows the effective change set"`
}

// SearchResponse is one page of hits, best first.
type SearchResponse struct {
	Hits       []SearchHit `json:"hits" desc:"The hits, best first, equal scores by ID"`
	NextCursor string      `json:"next_cursor,omitempty" desc:"The cursor of the next page; absent when no more remain"`
	View       string      `json:"view" desc:"The view it was read in: as-is, to-be or what-if:<change set ID>"`

	// command is the search command line of this page, without --cursor.
	command string
}

// SearchHit is one node found.
type SearchHit struct {
	ID      id.ID          `json:"id" desc:"The node's ID"`
	Kind    op.Kind        `json:"kind" desc:"The node's kind"`
	Handle  string         `json:"handle" desc:"The node's handle"`
	Path    []ShowPathStep `json:"path" desc:"Its ancestors in the tree, the root first; empty for a root or a kind outside the tree"`
	Field   string         `json:"field" desc:"Where the text matched best: \"handle\" or a field name"`
	Snippet string         `json:"snippet" desc:"The matched text, shortened around the match"`
	Score   int            `json:"score" desc:"The rank: higher is better"`
}

// searchMatch is the matcher's verdict on one node.
type searchMatch struct {
	score   int
	field   string
	snippet string
}

// runSearch checks the request, ranks every match, then cuts the page.
func runSearch(ctx *Ctx, req any) (any, error) {
	q := req.(SearchRequest)
	terms := index.Tokens(q.Text)
	if len(terms) == 0 {
		return nil, Usagef("nothing to search for: give some words").WithNext(searchDescribeNext)
	}
	var kinds []op.Kind
	for _, k := range q.Kinds {
		if !slices.Contains(op.Kinds(), op.Kind(k)) {
			return nil, Usagef("unknown kind %q: want one of %s", k, joinKinds()).WithNext(searchDescribeNext)
		}
		if !slices.Contains(kinds, op.Kind(k)) {
			kinds = append(kinds, op.Kind(k))
		}
	}
	limit := q.Limit
	if limit == 0 {
		limit = searchDefaultLimit
	}
	if limit < 0 || limit > searchMaxLimit {
		return nil, Usagef("limit %d is out of range: give 1 to %d", q.Limit, searchMaxLimit).WithNext(searchDescribeNext)
	}
	var afterScore int
	var afterID id.ID
	hasAfter := false
	if q.Cursor != "" {
		s, i, err := searchDecodeCursor(q.Cursor)
		if err != nil {
			return nil, err
		}
		afterScore, afterID, hasAfter = s, i, true
	}
	r, err := openReader(ctx, "search", q.View, defaultView)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	cands, err := r.Search(terms, kinds)
	if err != nil {
		return nil, err
	}
	type scored struct {
		n gmodel.Node
		m searchMatch
	}
	var all []scored
	for _, n := range cands {
		if m, ok := matchNode(n, terms); ok {
			all = append(all, scored{n, m})
		}
	}
	slices.SortFunc(all, func(a, b scored) int {
		if a.m.score != b.m.score {
			return b.m.score - a.m.score
		}
		return id.Compare(a.n.ID, b.n.ID)
	})
	if hasAfter {
		i := 0
		for i < len(all) && (all[i].m.score > afterScore || all[i].m.score == afterScore && id.Compare(all[i].n.ID, afterID) <= 0) {
			i++
		}
		all = all[i:]
	}
	resp := SearchResponse{Hits: []SearchHit{}, View: r.View(), command: searchCommand(q)}
	if len(all) > limit {
		all = all[:limit]
		last := all[limit-1]
		resp.NextCursor = searchEncodeCursor(last.m.score, last.n.ID)
	}
	for _, s := range all {
		as, err := r.Ancestors(s.n.ID)
		if err != nil {
			return nil, err
		}
		h := SearchHit{ID: s.n.ID, Kind: s.n.Kind, Handle: s.n.Handle, Path: []ShowPathStep{},
			Field: s.m.field, Snippet: s.m.snippet, Score: s.m.score}
		for _, a := range as {
			h.Path = append(h.Path, ShowPathStep{ID: a.ID, Handle: a.Handle})
		}
		resp.Hits = append(resp.Hits, h)
	}
	return resp, nil
}

// matchNode scores n against terms: every term must be a word or the start
// of a word of the handle or of some field. Each term scores its best
// place; the node's score is the sum, and its matched field the place that
// gave the most, the handle first on a tie, then field names in order.
func matchNode(n gmodel.Node, terms []string) (searchMatch, bool) {
	names := index.FieldNames(n)
	texts := make([]string, 0, len(names)+1)
	where := make([]string, 0, len(names)+1)
	texts, where = append(texts, n.Handle), append(where, "handle")
	for _, k := range names {
		texts, where = append(texts, index.FieldText(n.Fields[k])), append(where, k)
	}
	toks := make([][]string, len(texts))
	for i, t := range texts {
		toks[i] = index.Tokens(t)
	}
	per := make([]int, len(texts))
	total := 0
	for _, term := range terms {
		best := 0
		for i := range texts {
			s := termScore(toks[i], term, i == 0)
			per[i] += s
			best = max(best, s)
		}
		if best == 0 {
			return searchMatch{}, false
		}
		total += best
	}
	at := 0
	for i := range per {
		if per[i] > per[at] {
			at = i
		}
	}
	return searchMatch{score: total, field: where[at], snippet: snippetOf(texts[at], terms)}, true
}

// termScore is the best score of term in the words toks, 0 for none.
func termScore(toks []string, term string, handle bool) int {
	word, prefix := scoreFieldWord, scoreFieldPrefix
	if handle {
		word, prefix = scoreHandleWord, scoreHandlePrefix
	}
	best := 0
	for _, t := range toks {
		switch {
		case t == term:
			return word
		case strings.HasPrefix(t, term):
			best = prefix
		}
	}
	return best
}

// The snippet's window, in characters either side of the match.
const (
	snippetBefore = 30
	snippetAfter  = 60
)

// snippetOf shortens text to a window around the first of terms it holds,
// whitespace collapsed, with "..." where it is cut.
func snippetOf(text string, terms []string) string {
	text = strings.Join(strings.Fields(text), " ")
	rs := []rune(text)
	low := []rune(strings.ToLower(text))
	start := 0
	if len(low) == len(rs) {
		for _, t := range terms {
			if i := strings.Index(string(low), t); i >= 0 {
				start = len([]rune(string(low)[:i]))
				break
			}
		}
	}
	from := max(0, start-snippetBefore)
	to := min(len(rs), start+snippetAfter)
	out := string(rs[from:to])
	if from > 0 {
		out = "..." + out
	}
	if to < len(rs) {
		out += "..."
	}
	return out
}

// searchEncodeCursor is the cursor that goes on after the hit with score
// and ID i.
func searchEncodeCursor(score int, i id.ID) string {
	return base64.RawURLEncoding.EncodeToString([]byte(strconv.Itoa(score) + ":" + i.String()))
}

// searchDecodeCursor reverses searchEncodeCursor; anything that is not a
// cursor is a usage error.
func searchDecodeCursor(c string) (int, id.ID, error) {
	if b, err := base64.RawURLEncoding.DecodeString(c); err == nil {
		if s, rest, ok := strings.Cut(string(b), ":"); ok {
			if score, err := strconv.Atoi(s); err == nil {
				if i, err := id.Parse(rest); err == nil {
					return score, i, nil
				}
			}
		}
	}
	return 0, id.ID{}, Usagef("bad cursor %q: give a next_cursor from search", c).WithNext(searchDescribeNext)
}

// searchCommand is the command line that asks for q again, without its
// cursor.
func searchCommand(q SearchRequest) string {
	parts := []string{"rhiza", "search", q.Text}
	for _, k := range q.Kinds {
		parts = append(parts, "--kind", k)
	}
	if q.Limit != 0 {
		parts = append(parts, "--limit", strconv.Itoa(q.Limit))
	}
	if q.View != "" {
		parts = append(parts, "--view", q.View)
	}
	return commandLine(parts...)
}

// Human lists the hits, each as a line and its snippet, then the command
// for the next page when more remain.
func (r SearchResponse) Human(w io.Writer) {
	if len(r.Hits) == 0 {
		fmt.Fprintln(w, "No nodes match.")
		return
	}
	for _, h := range r.Hits {
		fmt.Fprintf(w, "- %s %s %s\n", h.Kind, h.Handle, h.ID)
		if len(h.Path) > 0 {
			hs := make([]string, len(h.Path))
			for i, p := range h.Path {
				hs[i] = p.Handle
			}
			fmt.Fprintf(w, "  path: %s\n", strings.Join(hs, " / "))
		}
		fmt.Fprintf(w, "  %s: %s\n", h.Field, h.Snippet)
	}
	if r.NextCursor != "" {
		cmd := r.command
		if cmd == "" {
			cmd = "rhiza search"
		}
		fmt.Fprintf(w, "\nnext: %s --cursor %s\n", cmd, r.NextCursor)
	}
}

// Search returns the live nodes that may match terms (lower-case words),
// of kinds when given, sorted by ID. The stored views read the FTS table;
// what-if reads the overlay.
func (r *reader) Search(terms []string, kinds []op.Kind) ([]gmodel.Node, error) {
	if r.ov != nil {
		ns, err := r.ov.Search(terms, kinds)
		return ns, readError(err)
	}
	ns, err := r.x.Search(r.view.stored, terms, kinds)
	return ns, readError(err)
}
