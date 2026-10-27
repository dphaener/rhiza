// Package guide holds the agent guide that ships in the binary: the rhiza
// skill, the rhiza-scrub skill, the rhiza-plan skill and the topics an
// agent reads on demand.
//
// The guide comes in three layers, each with a size cap measured in
// estimated tokens (Tokens): the session start a hook prints (StartCap),
// a skill (SkillCap) and a topic (TopicCap). The texts are markdown files
// embedded from text/; this package only looks them up and measures them.
package guide

import (
	"embed"
	"slices"
)

// The size caps, in estimated tokens.
const (
	// SkillCap is the most a skill may cost.
	SkillCap = 2000
	// TopicCap is the most a topic may cost.
	TopicCap = 800
	// StartCap is the most the session-start context may cost.
	StartCap = 100
)

// MaxDescription is the longest a skill or topic description may be, in
// bytes: the limit on a SKILL.md description.
const MaxDescription = 1024

// Text is one text of the guide.
type Text struct {
	// Name is what rhiza guide takes: rhiza, scrub, plan or a topic name.
	Name string
	// Description says when to load it; for a skill it is the SKILL.md
	// description.
	Description string
	// Body is the markdown.
	Body string
	// Skill is true for the skills, false for topics.
	Skill bool
}

// Cap returns the size cap that applies to t, in estimated tokens.
func (t Text) Cap() int {
	if t.Skill {
		return SkillCap
	}
	return TopicCap
}

//go:embed text/*.md
var files embed.FS

// entry is one text's name, file and description.
type entry struct {
	name, file, description string
	skill                   bool
}

// skills lists the skills, the main one first.
var skills = []entry{
	{
		name: "rhiza",
		file: "rhiza.md",
		description: "Work with the product model of a repository that has a .rhiza/ folder: " +
			"capabilities, intents, constraints, decisions, questions, evidence and sources, written " +
			"through change sets with the rhiza CLI. Load it before changing code in such a repo, " +
			"when asked to add or change a capability, feature, requirement or spec, record a decision " +
			"or ADR, raise a question, record what existing code does, or work on a change set.",
		skill: true,
	},
	{
		name: "scrub",
		file: "scrub.md",
		description: "Turn a raw conversation (a chat transcript, meeting notes, an interview) into a " +
			"Rhiza Source: summarise it, keep the owner's exact words wherever something was decided, " +
			"remove secrets, personal data and noise, then store it with rhiza source add. Load it " +
			"when asked to record, scrub or save a conversation as a source in a repo with a .rhiza/ folder.",
		skill: true,
	},
	{
		name: "plan",
		file: "plan.md",
		description: "Plan with the owner straight into a Rhiza product model: a milestone (an intent, " +
			"its done points and the change sets that serve it) or one change set's leaves, each with " +
			"an acceptance check, Decisions quoting the owner and parked questions, all written with " +
			"rhiza apply, the interview kept as a Source. Load it when you run as the planner of a plan " +
			"run, or are asked to plan a milestone or a change set's leaves with the owner in a repo " +
			"with a .rhiza/ folder.",
		skill: true,
	},
}

// topics lists the topics, sorted by name.
var topics = []entry{
	{name: "apply", file: "apply.md", description: "The rhiza apply file format: many writes in one JSON file."},
	{name: "arbiter", file: "arbiter.md", description: "The arbiter: deciding one question in an autonomous run, the one answering Decision and its wording and reasoning, then done."},
	{name: "build", file: "build.md", description: "Building a change set's leaves: a builder's one commit with a Leaf: trailer, outside work as a question or a discovered-from draft change set, a verifier's Evidence or findings. Load it when you run as a builder or verifier of a build run."},
	{name: "carving", file: "carving.md", description: "Carving: moving nodes from one change set to another under the same IDs, what moves, the rules."},
	{name: "change-sets", file: "change-sets.md", description: "Change sets: intention and record, the lifecycle, which change set a write goes into, --record."},
	{name: "plan-review", file: "plan-review.md", description: "Plan review: finding a plan's change sets, recording findings of five kinds, a critic's verdicts and a reviser's fixes."},
	{name: "reopening", file: "reopening.md", description: "Reopening: taking a change set landed only on its branch back to accepted, the refusal once it is on main, never deleting log entries."},
	{name: "review", file: "review.md", description: "Reviewing a built change set: five checks' findings with severities from the merge-base, a critic's verdicts that add no finding, fixers' commits with a Leaf: or Review: trailer, outside work as a question or a discovered-from draft change set, the QA check's Evidence. Load it when you run as a check, critic or fixer of a review run."},
	{name: "sources", file: "sources.md", description: "Source nodes: the stored fields, the raw material's hash, and how scrub feeds rhiza source add."},
}

// texts holds every text in List order, read once from the embedded files.
var texts = func() []Text {
	var out []Text
	for _, e := range slices.Concat(skills, topics) {
		b, err := files.ReadFile("text/" + e.file)
		if err != nil {
			panic("guide: " + err.Error())
		}
		out = append(out, Text{Name: e.name, Description: e.description, Body: string(b), Skill: e.skill})
	}
	return out
}()

// List returns every text: the rhiza, scrub and plan skills, then the
// topics sorted by name.
func List() []Text {
	return slices.Clone(texts)
}

// Lookup returns the text called name. "" and "rhiza" name the rhiza
// skill, "scrub" the rhiza-scrub skill, "plan" the rhiza-plan skill;
// anything else is a topic name.
func Lookup(name string) (Text, bool) {
	if name == "" {
		name = "rhiza"
	}
	for _, t := range texts {
		if t.Name == name {
			return t, true
		}
	}
	return Text{}, false
}

// Tokens estimates how many tokens s costs: its length in bytes over four,
// rounded up.
func Tokens(s string) int {
	return (len(s) + 3) / 4
}
