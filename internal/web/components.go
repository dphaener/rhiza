package web

import "html/template"

// The data the shared partials (templates/partials.html) are rendered with.

// Pill is a status pill. State picks the colour and glyph: for runs running,
// waiting, done, failed or interrupted; for change sets draft, proposed,
// accepted, landed or withdrawn. Text is the word beside the glyph.
type Pill struct{ State, Text string }

// Panel is a titled panel with a kicker. Body is trusted markup, normally
// the output of another partial.
type Panel struct {
	Kicker, Title, Aside string
	Flush, Waiting       bool
	Body                 template.HTML
}

// Glyph is a node's typed glyph: intent, capability, decision, constraint or
// question.
type Glyph struct {
	Kind  string
	Small bool
}

// Row is a dense row: a title linking to Href (the whole row follows it),
// an optional sub line, a status pill, the current step and when it changed.
// It shows no ID (constraint owner-surfaces-no-ids): Href may carry one, the
// text never does. Datetime, if set, makes Time a live relative-time element.
type Row struct {
	Href, Title, Sub string
	Pill             Pill
	Step             string
	Live             bool
	Time, Datetime   string
}

// Step is one step of a run's timeline. State is done, running, waiting,
// failed, pending or skipped.
type Step struct {
	Name, State string
	Pill        Pill
	Duration    string
	Note        string
	Sessions    []Session
}

// Session is an agent session under a step, linked by who ran it, never
// by its ID (constraint owner-surfaces-no-ids).
type Session struct{ Href, Who, Model, Usage string }

// Term is the always-dark terminal block. Live adds the cursor and the live
// dot.
type Term struct {
	Label, Head, Path string
	Live              bool
	Lines             []TermLine
}

// TermLine is one line of terminal output. Kind colours Text: k (path),
// w (warning) or g (success); empty is plain.
type TermLine struct{ Time, Kind, Text string }

// Empty is the empty state: a title and a line saying what to do.
type Empty struct{ Title, Text string }
