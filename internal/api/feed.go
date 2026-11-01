package api

// This file holds feed: the change sets that landed, oldest land first, a
// page at a time. A change set sits at its latest land entry, so one that
// was reopened and landed again appears once, at its new land. The cursor
// is opaque to callers: it is the base64url (unpadded) form of the land
// entry ID the page before ended at.

import (
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
)

func init() {
	register(Operation{
		Name:     "feed",
		Summary:  "List the change sets that landed, in land order, a page at a time",
		Request:  reflect.TypeFor[FeedRequest](),
		Response: reflect.TypeFor[FeedResponse](),
		Run:      runFeed,
	})
}

// The page sizes of feed.
const (
	feedDefaultLimit = 100
	feedMaxLimit     = 1000
)

// feedDescribeNext is the next command for a bad feed request.
const feedDescribeNext = "rhiza describe feed"

// feedPRNumber picks the pull request number off the end of a pr field.
var feedPRNumber = regexp.MustCompile(`(\d+)/?$`)

// FeedRequest pages and filters the landed change sets.
type FeedRequest struct {
	Limit  int    `json:"limit,omitempty" desc:"The most change sets to return, 1 to 1000; default 100"`
	Cursor string `json:"cursor,omitempty" desc:"Where to go on: the next_cursor of the page before"`
	Since  string `json:"since,omitempty" desc:"Keep the change sets landed at or after this time, RFC 3339"`
}

// FeedResponse is one page of the feed.
type FeedResponse struct {
	ChangeSets []FeedItem `json:"change_sets" desc:"The landed change sets, oldest land first"`
	NextCursor string     `json:"next_cursor,omitempty" desc:"The cursor of the next page; absent when no more remain"`

	// command is the feed command line of this page, without --cursor, for
	// the next line.
	command string
}

// FeedItem is one landed change set.
type FeedItem struct {
	ID        id.ID          `json:"id" desc:"The change set's ID"`
	Title     string         `json:"title" desc:"The change set's title, else its intent's first clause"`
	Kind      changeset.Kind `json:"kind" desc:"intention or record"`
	Leaves    int            `json:"leaves" desc:"How many leaf capabilities it created"`
	LandedAt  string         `json:"landed_at" desc:"When it landed, RFC 3339 UTC"`
	Actor     string         `json:"actor" desc:"Who landed it: the actor's handle, else its ID"`
	PR        int            `json:"pr,omitempty" desc:"The pull request number from its PR Evidence; absent when there is none"`
	LandEntry id.ID          `json:"land_entry" desc:"The land entry the change set sits at"`
}

// runFeed folds the log, keeps the landed change sets at their latest land
// entry, orders them by it and cuts the page.
func runFeed(ctx *Ctx, req any) (any, error) {
	q := req.(FeedRequest)
	limit := q.Limit
	if limit == 0 {
		limit = feedDefaultLimit
	}
	if limit < 0 || limit > feedMaxLimit {
		return nil, Usagef("limit %d is out of range: give 1 to %d", q.Limit, feedMaxLimit).WithNext(feedDescribeNext)
	}
	var after id.ID
	if q.Cursor != "" {
		c, err := feedDecodeCursor(q.Cursor)
		if err != nil {
			return nil, err
		}
		after = c
	}
	var since time.Time
	if q.Since != "" {
		t, err := time.Parse(time.RFC3339, q.Since)
		if err != nil {
			return nil, Usagef("bad since %q: give an RFC 3339 time such as 2026-01-02T15:04:05Z", q.Since).WithNext(feedDescribeNext)
		}
		since = t
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	m := w.model(id.ID{})
	var items []FeedItem
	for _, c := range w.changeSets() {
		if c.State != changeset.StateLanded {
			continue
		}
		land, ok := latestLand(c)
		if !ok || (!since.IsZero() && land.At.Before(since)) || id.Compare(land.Entry, after) <= 0 {
			continue
		}
		items = append(items, FeedItem{
			ID: c.ID, Title: c.Open.DisplayTitle(), Kind: c.Open.Kind,
			Leaves:    len(changeLeaves(m, c.Ops)),
			LandedAt:  rfc3339(land.At),
			Actor:     feedActor(m, land.Actor),
			PR:        feedPR(m, c, land.Entry),
			LandEntry: land.Entry,
		})
	}
	slices.SortFunc(items, func(a, b FeedItem) int { return id.Compare(a.LandEntry, b.LandEntry) })
	resp := FeedResponse{ChangeSets: []FeedItem{}, command: feedCommand(q)}
	if len(items) > limit {
		items = items[:limit]
		resp.NextCursor = feedEncodeCursor(items[limit-1].LandEntry)
	}
	resp.ChangeSets = append(resp.ChangeSets, items...)
	return resp, nil
}

// latestLand is the last land move that applied to c.
func latestLand(c changeset.ChangeSet) (changeset.Move, bool) {
	for _, mv := range slices.Backward(c.Lifecycle) {
		if mv.Type == entry.TypeLand {
			return mv, true
		}
	}
	return changeset.Move{}, false
}

// feedActor is the handle of actor a, else its ID.
func feedActor(m Model, a id.ID) string {
	if n, ok := m.Node(a); ok && n.Handle != "" {
		return n.Handle
	}
	return a.String()
}

// feedPR is the pull request number of c's PR Evidence, 0 when there is
// none. Evidence c's own ops create comes first; else Evidence recorded for
// the capabilities c changed, the first written after c landed (a record
// change set written once it had landed).
func feedPR(m Model, c changeset.ChangeSet, land id.ID) int {
	for _, o := range slices.Backward(c.Ops) {
		if o.Type != op.TypeNodeCreate || o.Kind != op.KindEvidence {
			continue
		}
		if n, ok := m.Node(o.Node); ok && !n.Retired {
			if pr := prNumber(n); pr > 0 {
				return pr
			}
		}
	}
	var best Node
	for _, capID := range changedCapabilities(m, c.Ops) {
		for _, l := range m.LinksTo(capID, op.LinkVerifies) {
			n, ok := m.Node(l.From)
			if !ok || n.Kind != op.KindEvidence || n.Retired || id.Compare(n.ID, land) <= 0 || prNumber(n) == 0 {
				continue
			}
			if best.ID.IsZero() || id.Compare(n.ID, best.ID) < 0 {
				best = n
			}
		}
	}
	if best.ID.IsZero() {
		return 0
	}
	return prNumber(best)
}

// prNumber is the pull request number in node n's pr field, 0 if none.
func prNumber(n Node) int {
	raw, ok := n.Fields["pr"]
	if !ok {
		return 0
	}
	g := feedPRNumber.FindStringSubmatch(strings.TrimSpace(fieldText(raw)))
	if g == nil {
		return 0
	}
	v, err := strconv.Atoi(g[1])
	if err != nil {
		return 0
	}
	return v
}

// feedEncodeCursor is the cursor that goes on after land entry e.
func feedEncodeCursor(e id.ID) string {
	return base64.RawURLEncoding.EncodeToString([]byte(e.String()))
}

// feedDecodeCursor is the land entry cursor c goes on after; anything that
// is not a cursor is a usage error.
func feedDecodeCursor(c string) (id.ID, error) {
	b, err := base64.RawURLEncoding.DecodeString(c)
	if err == nil {
		var i id.ID
		if i, err = id.Parse(string(b)); err == nil {
			return i, nil
		}
	}
	return id.ID{}, Usagef("bad cursor %q: give a next_cursor from feed", c).WithNext(feedDescribeNext)
}

// feedCommand is the command line that asks for q again, without its cursor.
func feedCommand(q FeedRequest) string {
	parts := []string{"rhiza", "feed"}
	if q.Limit != 0 {
		parts = append(parts, "--limit", strconv.Itoa(q.Limit))
	}
	if q.Since != "" {
		parts = append(parts, "--since", q.Since)
	}
	return commandLine(parts...)
}

// Human lists the page one change set a line, then the command for the next
// page when more remain.
func (r FeedResponse) Human(w io.Writer) {
	if len(r.ChangeSets) == 0 {
		fmt.Fprintln(w, "Nothing landed.")
		return
	}
	for _, c := range r.ChangeSets {
		line := fmt.Sprintf("- %s %s %s %s: %s (%d %s)", c.LandedAt, c.ID, c.Kind, c.Actor, oneLine(c.Title), c.Leaves, feedPlural(c.Leaves))
		if c.PR > 0 {
			line += fmt.Sprintf(" #%d", c.PR)
		}
		fmt.Fprintln(w, line)
	}
	if r.NextCursor != "" {
		cmd := r.command
		if cmd == "" {
			cmd = "rhiza feed"
		}
		fmt.Fprintf(w, "\nnext: %s --cursor %s\n", cmd, r.NextCursor)
	}
}

func feedPlural(n int) string {
	if n == 1 {
		return "leaf"
	}
	return "leaves"
}
