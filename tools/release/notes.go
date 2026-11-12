package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Notes drafts the release notes from the change sets landed after since,
// grouped by the milestone each serves. Breaking changes (marked, with their
// reasons, and the detected differences) and bug fixes (severe ones flagged)
// are called out first. The output depends only on its inputs.
func Notes(h History, since time.Time, marks map[string]Marks, detected []string) string {
	var ids []string
	for id, l := range h.Landed {
		if l.At.After(since) {
			ids = append(ids, id)
		}
	}
	sort.Slice(ids, func(i, j int) bool {
		a, b := h.Landed[ids[i]], h.Landed[ids[j]]
		if !a.At.Equal(b.At) {
			return a.At.Before(b.At)
		}
		return a.ID < b.ID
	})
	title := func(id string) string {
		if t := strings.TrimSpace(h.Landed[id].Intent); t != "" {
			return t
		}
		return id
	}
	milestoneOf := map[string]string{}
	var handles []string
	for hd := range h.Milestones {
		handles = append(handles, hd)
	}
	sort.Strings(handles)
	for _, hd := range handles {
		for _, cs := range h.Milestones[hd] {
			if _, ok := milestoneOf[cs]; !ok {
				milestoneOf[cs] = hd
			}
		}
	}
	var b strings.Builder
	section := func(head string, lines []string) {
		if len(lines) == 0 {
			return
		}
		fmt.Fprintf(&b, "## %s\n\n", head)
		for _, l := range lines {
			fmt.Fprintf(&b, "- %s\n", l)
		}
		b.WriteString("\n")
	}
	var breaking, fixes []string
	byMilestone := map[string][]string{}
	var other []string
	for _, id := range ids {
		m := marks[id]
		for _, r := range m.Breaking {
			breaking = append(breaking, fmt.Sprintf("%s: %s", title(id), r))
		}
		if h.Landed[id].RequestKind == "bug" {
			line := title(id)
			if m.IsSevere() {
				line += " (severe: " + strings.Join(m.Severe, "; ") + ")"
			}
			fixes = append(fixes, line)
			continue
		}
		if hd, ok := milestoneOf[id]; ok {
			byMilestone[hd] = append(byMilestone[hd], title(id))
		} else {
			other = append(other, title(id))
		}
	}
	breaking = append(breaking, detected...)
	section("Breaking changes", breaking)
	for _, hd := range handles {
		if len(byMilestone[hd]) == 0 {
			continue
		}
		head := "Milestone " + hd
		if s := h.Statements[hd]; s != "" {
			head += ": " + s
		}
		section(head, byMilestone[hd])
	}
	section("Other changes", other)
	section("Bug fixes", fixes)
	if b.Len() == 0 {
		return "No changes.\n"
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}
