package api

// This file holds feed: the change sets that landed, oldest land first, a
// page at a time. A change set sits at its latest land entry, so one that
// was reopened and landed again appears once, at its new land. The cursor
// is opaque to callers: it is the base64url (unpadded) form of the land
// entry ID the page before ended at.

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/canon"
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
	All    bool   `json:"all,omitempty" flag:"all" desc:"Include the change sets whose ops only create or update actor nodes, which are left out by default"`
	Follow bool   `json:"follow,omitempty" flag:"follow" desc:"Stream each change set as it lands until interrupted: with neither cursor nor since only those landing after the start, else the backlog after that point first; one FeedItem per line with --json; limit then stops after that many, default unlimited"`
}

// FeedResponse is one page of the feed.
type FeedResponse struct {
	ChangeSets []FeedItem `json:"change_sets" desc:"The landed change sets, oldest land first"`
	NextCursor string     `json:"next_cursor,omitempty" desc:"The cursor of the next page; absent when no more remain"`

	// streamed is set when the page went to Ctx.Out line by line (follow):
	// callers print nothing more.
	streamed bool

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
	if q.Follow {
		if limit < 0 {
			return nil, Usagef("limit %d is out of range: give 1 or more", q.Limit).WithNext(feedDescribeNext)
		}
	} else if limit < 0 || limit > feedMaxLimit {
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
	if q.Follow {
		return runFeedFollow(ctx, q, after, since)
	}
	items, err := feedItems(ctx, q.All, since, after)
	if err != nil {
		return nil, err
	}
	resp := FeedResponse{ChangeSets: []FeedItem{}, command: feedCommand(q)}
	if len(items) > limit {
		items = items[:limit]
		resp.NextCursor = feedEncodeCursor(items[limit-1].LandEntry)
	}
	resp.ChangeSets = append(resp.ChangeSets, items...)
	return resp, nil
}

// feedItems folds the log and returns the landed change sets at their latest
// land entry, past land entry after and not landed before since (zero for
// either: no bound), ordered by land entry.
func feedItems(ctx *Ctx, all bool, since time.Time, after id.ID) ([]FeedItem, error) {
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
		if !all && feedActorOnly(m, c) {
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
	return items, nil
}

// feedActorOnly reports whether every op of c creates or updates an actor
// node: the record of an agent or person joining, not work.
func feedActorOnly(m Model, c changeset.ChangeSet) bool {
	if len(c.Ops) == 0 {
		return false
	}
	for _, o := range c.Ops {
		if o.Type != op.TypeNodeCreate && o.Type != op.TypeNodeUpdate {
			return false
		}
		if o.Type == op.TypeNodeCreate && o.Kind == op.KindActor {
			continue
		}
		if n, ok := m.Node(o.Node); !ok || n.Kind != op.KindActor {
			return false
		}
	}
	return true
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
	if q.All {
		parts = append(parts, "--all")
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
		fmt.Fprintln(w, feedLine(c))
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

// feedLine is the line of change set c in human output.
func feedLine(c FeedItem) string {
	line := fmt.Sprintf("- %s %s %s %s: %s (%d %s)", c.LandedAt, c.ID, c.Kind, c.Actor, oneLine(c.Title), c.Leaves, feedPlural(c.Leaves))
	if c.PR > 0 {
		line += fmt.Sprintf(" #%d", c.PR)
	}
	return line
}

// feedPoll is the interval at which follow refolds the log.
const feedPoll = 2 * time.Second

// feedTicker ticks once per follow poll; tests replace it.
var feedTicker = func() (<-chan time.Time, func()) {
	t := time.NewTicker(feedPoll)
	return t.C, t.Stop
}

// Streamed reports that the response went to Ctx.Out as it happened, so a
// caller prints nothing for it.
func (r FeedResponse) Streamed() bool { return r.streamed }

// runFeedFollow streams the change sets landing past the last one emitted
// to ctx.Out, one line each (a FeedItem as one line of JSON when ctx.JSON,
// else feed's '- ' line), at once and then on each tick, until interrupted
// or q.Limit change sets are out. With neither cursor nor since it starts at
// the current end of the feed.
func runFeedFollow(ctx *Ctx, q FeedRequest, after id.ID, since time.Time) (any, error) {
	if ctx.Out == nil {
		return nil, Usagef("feed --follow needs somewhere to stream to").WithNext(feedDescribeNext)
	}
	if q.Cursor == "" && q.Since == "" {
		// The end of the feed counts every landed change set, hidden or not.
		items, err := feedItems(ctx, true, time.Time{}, id.ID{})
		if err != nil {
			return nil, err
		}
		if len(items) > 0 {
			after = items[len(items)-1].LandEntry
		}
	}
	var interrupted <-chan struct{}
	if ctx.Interrupts != nil {
		ch, stop := ctx.Interrupts()
		defer stop()
		interrupted = ch
	}
	ticks, stopTicks := feedTicker()
	defer stopTicks()
	done := FeedResponse{ChangeSets: []FeedItem{}, streamed: true}
	emitted := 0
	for {
		items, err := feedItems(ctx, q.All, since, after)
		if err != nil {
			return nil, err
		}
		for _, c := range items {
			var line []byte
			if ctx.JSON {
				b, err := canon.Marshal(c)
				if err != nil {
					return nil, Errorf("encode feed item: %v", err)
				}
				var buf bytes.Buffer
				if err := json.Compact(&buf, bytes.TrimSpace(b)); err != nil {
					return nil, Errorf("encode feed item: %v", err)
				}
				buf.WriteByte('\n')
				line = buf.Bytes()
			} else {
				line = []byte(feedLine(c) + "\n")
			}
			if _, err := ctx.Out.Write(line); err != nil {
				return nil, Errorf("write feed: %v", err)
			}
			after = c.LandEntry
			emitted++
			if q.Limit > 0 && emitted >= q.Limit {
				return done, nil
			}
		}
		select {
		case <-interrupted:
			return done, nil
		case <-ticks:
		}
	}
}
