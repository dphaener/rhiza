package api

// This file holds the repair operation: it restores the damaged log entries
// render refuses (see findDamaged). The authority for an entry tracked at
// HEAD is HEAD's blob; for any other the worktree's persistent index, which
// keeps the canonical bytes of each entry it read well, whatever happened to
// the file since. Nothing is deleted, and every file is written through a
// temp file and a rename at the entry mode.

import (
	"fmt"
	"io"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/store"
)

// RepairRequest takes nothing: repair works in the caller's root; --dry-run
// names what it would restore.
type RepairRequest struct{}

// RepairedEntry is one damaged entry and what repair did, or would do, with it.
type RepairedEntry struct {
	Path   string `json:"path" desc:"The entry file, relative to the project root"`
	Reason string `json:"reason" desc:"What was wrong with it"`
	From   string `json:"from,omitempty" desc:"Where the original bytes come from: HEAD or index; empty when there are none"`
	Why    string `json:"why,omitempty" desc:"Why it cannot be restored; empty when it can"`
}

// RepairResponse lists the damaged entries.
type RepairResponse struct {
	DryRun   bool            `json:"dry_run" desc:"True when nothing was written"`
	Restored []RepairedEntry `json:"restored" desc:"The entries restored, or with --dry-run the ones that would be"`
	Lost     []RepairedEntry `json:"lost" desc:"The entries that cannot be restored, with why"`
}

// Human names each entry repaired, or with --dry-run each one repair would
// restore, and where its bytes come from.
func (r RepairResponse) Human(w io.Writer) {
	if len(r.Restored) == 0 && len(r.Lost) == 0 {
		fmt.Fprintln(w, "No damaged log entries; nothing to repair.")
		return
	}
	for _, e := range r.Restored {
		fmt.Fprintf(w, "%s: %s from %s (%s)\n", e.Path, restoredVerb(r.DryRun), e.From, e.Reason)
	}
	for _, e := range r.Lost {
		fmt.Fprintf(w, "%s: not restorable: %s\n", e.Path, e.Why)
	}
	if r.DryRun {
		fmt.Fprintln(w, "Dry run: nothing was written.")
	}
}

func init() {
	register(Operation{
		Name:     "repair",
		Summary:  "Restore the damaged log entries render refuses, from HEAD or the persistent index",
		Request:  reflect.TypeFor[RepairRequest](),
		Response: reflect.TypeFor[RepairResponse](),
		Run:      runRepair,
	})
}

func runRepair(ctx *Ctx, _ any) (any, error) {
	st, err := openStore(ctx.Root)
	if err != nil {
		return nil, err
	}
	resp := RepairResponse{DryRun: ctx.DryRun, Restored: []RepairedEntry{}, Lost: []RepairedEntry{}}
	ds := findDamaged(ctx, st.ReadLog())
	if len(ds) == 0 {
		return resp, nil
	}
	var x *index.Index
	defer func() {
		if x != nil {
			x.Close()
		}
	}()
	for _, d := range ds {
		e := RepairedEntry{Path: d.Path, Reason: d.Reason}
		data, from, why := originalBytes(ctx, d, &x)
		if why != "" {
			e.Why = why
			resp.Lost = append(resp.Lost, e)
			continue
		}
		e.From = from
		if !ctx.DryRun {
			dst := filepath.Join(ctx.Root, filepath.FromSlash(d.Path))
			if err := store.WriteEntryFile(dst, data); err != nil {
				e.From, e.Why = "", fmt.Sprintf("write failed: %v", err)
				resp.Lost = append(resp.Lost, e)
				continue
			}
		}
		resp.Restored = append(resp.Restored, e)
	}
	if len(resp.Lost) > 0 {
		var lines []string
		for _, e := range resp.Restored {
			lines = append(lines, fmt.Sprintf("%s: %s from %s", e.Path, restoredVerb(ctx.DryRun), e.From))
		}
		for _, e := range resp.Lost {
			lines = append(lines, fmt.Sprintf("%s: not restorable: %s", e.Path, e.Why))
		}
		return nil, Errorf("repair: %d entries cannot be restored:\n%s", len(resp.Lost), strings.Join(lines, "\n"))
	}
	return resp, nil
}

func restoredVerb(dry bool) string {
	if dry {
		return "would restore"
	}
	return "restored"
}

// originalBytes finds the bytes d's file should hold: HEAD's when the file
// is tracked there, else the index's, either only when they are canonical and
// name the file's own change set and ID (an entry can be committed already
// damaged); a file that is not an entry, such as repo.json, takes HEAD's
// unchecked. why says why there are none. The index is opened on first need
// into x.
func originalBytes(ctx *Ctx, d damaged, x **index.Index) (data []byte, from, why string) {
	cs, eid, ok := entryPathIDs(d.Path)
	headDamaged := false
	if inWorkTree(ctx.Git) {
		if out, err := ctx.Git("show", "HEAD:./"+d.Path); err == nil {
			if !ok || isEntryBytes([]byte(out), cs, eid) {
				return []byte(out), "HEAD", ""
			}
			headDamaged = true
		}
	}
	if !ok {
		return nil, "", "not a log/<change set>/<entry>.json path"
	}
	data, why = indexBytes(ctx, cs, eid, x)
	if why != "" {
		head := "not tracked at HEAD"
		if headDamaged {
			head = "HEAD's copy is damaged too"
		}
		return nil, "", head + " and " + why
	}
	return data, "index", ""
}

// indexBytes reads the canonical bytes of entry eid in change set cs from the
// persistent index, opened on first need into x. why says why there are none.
func indexBytes(ctx *Ctx, cs, eid id.ID, x **index.Index) (data []byte, why string) {
	if *x == nil {
		xi, err := index.Open(ctx.Root, ctx.IndexPath)
		if err != nil {
			return nil, "the index cannot be opened: " + err.Error()
		}
		*x = xi
	}
	raw, found, err := (*x).StoredEntry(cs, eid)
	if err != nil {
		return nil, "the index cannot be read: " + err.Error()
	}
	if !found {
		return nil, "the index never saw it"
	}
	if !isEntryBytes(raw, cs, eid) {
		return nil, "the index's stored bytes are not this entry's canonical bytes"
	}
	return raw, ""
}

// isEntryBytes reports whether data is the canonical encoding of entry eid in
// change set cs.
func isEntryBytes(data []byte, cs, eid id.ID) bool {
	e, err := entry.DecodeCanonical(data)
	return err == nil && e.ID == eid && e.ChangeSet == cs
}

// entryPathIDs parses .rhiza/log/<change set>/<entry>.json.
func entryPathIDs(p string) (cs, eid id.ID, ok bool) {
	rest, found := strings.CutPrefix(p, store.Dir+"/"+store.LogDir+"/")
	if !found {
		return cs, eid, false
	}
	dir, file := path.Split(rest)
	name, found := strings.CutSuffix(file, ".json")
	if !found {
		return cs, eid, false
	}
	cs, err1 := id.Parse(strings.TrimSuffix(dir, "/"))
	eid, err2 := id.Parse(name)
	return cs, eid, err1 == nil && err2 == nil
}
