package slice

// This file fits a slice to a token budget.
//
// A token is estimated as 4 bytes of the printed output, rounded up
// (Tokens); there is no tokenizer. The budget is measured on the bytes
// actually printed, markdown or JSON, footer included.
//
// Each section has a tier: full, brief (each item as its handle and short
// ID only) or count (no items; only its count, in the footer). Fit renders
// every section at full. While that is over budget it shrinks the sections
// one at a time to brief, in this order, re-measuring after each: In
// flight, Neighbours, Children, Evidence, Anchors, Decisions, Context. If
// that is still over it does the same pass down to count. Focus,
// Constraints and Open questions are never shrunk; an empty section is
// skipped. If the slice is still over budget with every other section at a
// count, it is printed anyway and marked over budget; that is not an error.
// Over budget is decided before the mark is printed: in JSON,
// "over_budget": true is a byte shorter than false, so an over-budget JSON
// slice may measure exactly the budget.
//
// Each cut section gets a command that shows what was left out. For
// Children in a node slice it is "rhiza tree <focus handle>". Anything
// else, Children in a change slice included (a change slice has no single
// tree to show), gets the slice command with " --budget <B>", B the uncut
// slice's tokens rounded up to the next 1,000 above it. The slice command
// is "rhiza slice <ref>" for a node slice and "rhiza slice --change
// <ref>" for a change slice, plus " --view <view>" when the view is not the
// request's default (DefaultView). A ref or handle that is not plain is
// single-quoted for the shell.

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/changeset"
)

// The views.
const (
	ViewAsIs   = "as-is"
	ViewToBe   = "to-be"
	ViewWhatIf = "what-if"
)

// DefaultView returns the view a request gets when it names none, from its
// effective change set: what-if of a draft or proposed one, to-be for an
// accepted one, and as-is for none or any other state (landed or terminal).
func DefaultView(cs *changeset.ChangeSet) string {
	if cs == nil {
		return ViewAsIs
	}
	switch cs.State {
	case changeset.StateDraft, changeset.StateProposed:
		return ViewWhatIf
	case changeset.StateAccepted:
		return ViewToBe
	}
	return ViewAsIs
}

// Tokens returns the estimated tokens of b: its length divided by 4,
// rounded up.
func Tokens(b []byte) int {
	return (len(b) + 3) / 4
}

// Tier is how much of a section is printed.
type Tier int

// The tiers, from most to least.
const (
	TierFull Tier = iota
	TierBrief
	TierCount
)

// String returns "full", "brief" or "count".
func (t Tier) String() string {
	switch t {
	case TierBrief:
		return "brief"
	case TierCount:
		return "count"
	}
	return "full"
}

// cutOrder is the order sections shrink in, at each tier.
var cutOrder = [...]part{secInFlight, secNeighbours, secChildren, secEvidence, secAnchors, secDecisions, secContext}

// Options says how to print a slice.
type Options struct {
	// Budget is the token budget.
	Budget int
	// JSON prints canonical JSON instead of markdown.
	JSON bool
	// DefaultView is the view the request would get without --view; the
	// footer's commands name the slice's view only when it differs.
	DefaultView string
	// Ref is the ref the commands and the JSON ref name: the focus handle
	// for a node slice or the change set's short ID for a change slice.
	// Empty, it is the focus handle or the change set's short ID.
	Ref string
}

// Omission is a section cut to fit the budget.
type Omission struct {
	Section string `json:"section" enum:"context,decisions,anchors,evidence,children,neighbours,in_flight" desc:"The section's key"`
	Count   int    `json:"count" desc:"How many items the section holds"`
	Tier    string `json:"tier" enum:"brief,count" desc:"brief: items as handle and short ID only; count: no items"`
	Next    string `json:"next" desc:"A command that shows what was left out"`
}

// Result is a slice as printed. Doc is set for JSON, and canon.Marshal of
// it is Bytes.
type Result struct {
	Bytes   []byte
	Tokens  int
	Over    bool
	Omitted []Omission
	Doc     *Doc
}

// Fit prints s within opts.Budget tokens, as the file comment says. It
// fails only when a field of a focus node is not valid JSON.
func Fit(s Slice, opts Options) (Result, error) {
	f := fitter{s: s, opts: opts, ref: opts.Ref}
	if f.ref == "" {
		switch {
		case s.ChangeSlice && s.Change != nil:
			f.ref = s.Change.Short
		case len(s.Focus) > 0:
			f.ref = s.Focus[0].Handle
		default:
			f.ref = s.Title
		}
	}
	var tiers [numSections]Tier
	uncut, err := f.render(tiers, false)
	if err != nil || uncut.Tokens <= opts.Budget {
		return uncut, err
	}
	f.next = (uncut.Tokens/1000 + 1) * 1000
	for _, t := range []Tier{TierBrief, TierCount} {
		for _, sec := range cutOrder {
			if s.count(sec) == 0 {
				continue
			}
			tiers[sec] = t
			res, err := f.render(tiers, false)
			if err != nil || res.Tokens <= opts.Budget {
				return res, err
			}
		}
	}
	return f.render(tiers, true)
}

type fitter struct {
	s    Slice
	opts Options
	ref  string
	next int // the --budget the commands suggest
}

// render prints the slice at tiers. The output states numbers about itself
// (the JSON tokens field, the footer's tokens over budget); render prints
// again until they match what it measures. Each number only grows from
// one pass to the next, and a bigger number is never fewer bytes, so this
// ends.
func (f *fitter) render(tiers [numSections]Tier, over bool) (Result, error) {
	lay := layout{tiers: tiers, budget: f.opts.Budget, over: over}
	for sec := range numSections {
		if tiers[sec] != TierFull {
			lay.omitted = append(lay.omitted, Omission{Section: sectionKeys[sec], Count: f.s.count(sec), Tier: tiers[sec].String(), Next: f.command(sec)})
		}
	}
	for {
		var b []byte
		var d *Doc
		if f.opts.JSON {
			d = doc(f.s, lay, f.s.View, f.ref)
			var err error
			if b, err = canon.Marshal(d); err != nil {
				return Result{}, fmt.Errorf("slice: encode: %w", err)
			}
		} else {
			b = markdown(f.s, lay)
		}
		n := Tokens(b)
		by := 0
		if over {
			by = n - f.opts.Budget
		}
		if by == lay.overBy && (n == lay.tokens || !f.opts.JSON) {
			return Result{Bytes: b, Tokens: n, Over: over, Omitted: lay.omitted, Doc: d}, nil
		}
		lay.tokens, lay.overBy = n, by
	}
}

// command returns the command that shows what cutting sec left out.
func (f *fitter) command(sec part) string {
	if sec == secChildren && !f.s.ChangeSlice && len(f.s.Focus) > 0 {
		return "rhiza tree " + shellQuote(f.s.Focus[0].Handle)
	}
	c := "rhiza slice " + shellQuote(f.ref)
	if f.s.ChangeSlice {
		c = "rhiza slice --change " + shellQuote(f.ref)
	}
	if f.s.View != f.opts.DefaultView {
		c += " --view " + shellQuote(f.s.View)
	}
	return fmt.Sprintf("%s --budget %d", c, f.next)
}

// plain matches a word the shell reads as itself.
var plain = regexp.MustCompile(`^[A-Za-z0-9_./:@%+=,-]+$`)

// shellQuote returns s as one shell word: itself when plain, else
// single-quoted.
func shellQuote(s string) string {
	if plain.MatchString(s) {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
