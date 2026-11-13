# Plan recipes

There is no plan file and no checkboxes: the model is the plan, written with `rhiza apply`.

## Scope by kind

Tell the scope from `kind` in `rhiza change show --json`:

- `intention`: one change set serving an existing intent. Plan its leaves, within its
  intent, into `RHIZA_CHANGE`; anything outside the intent becomes a question.
  If it serves no intent yet (`"serves": []` in `rhiza change show --json`, the change
  set's own serves links; slice's `serves` follows its nodes), settle that first, before any leaf:
  pick a landed intent it belongs under, or write a new one (`"kind":"intent"`,
  under its parent), with the owner. Then
  `rhiza link "$RHIZA_CHANGE" serves <intent> --change "$RHIZA_CHANGE"`.
  Link a milestone only when `RHIZA_MILESTONE` names one (the `serves` link must reach it,
  directly or through an intent under it) or an open milestone fits. Other work is built as
  backlog by `.rhiza/policy.toml`'s `[schedule]` `backlog` rule, unless it is "skip": then
  the `serves` link must reach a milestone, never an umbrella or area intent under no
  milestone.
- `record`: a draft opened for the plan: plan a milestone (`rhiza guide plan-interview`).

## Leaves

A leaf is a capability with `statement` (one sentence) and `acceptance`, linked
`decided-by` its Decision:

```json
{"writes":[
  {"op":"add","kind":"decision","handle":"csv-streaming","fields":[
    "statement=Stream CSV rows instead of building the file in memory",
    "quote=Stream it, reports get huge"]},
  {"op":"link","from":"csv-streaming","type":"derived-from","to":"<source>"},
  {"op":"add","kind":"capability","handle":"export-csv","parent":"reports","fields":[
    "statement=Users export any report as CSV",
    "acceptance=A test exports a 100k-row report with flat memory use"]},
  {"op":"link","from":"export-csv","type":"decided-by","to":"csv-streaming"}
]}
```

Write each acceptance as the tests covering that leaf, never the repo's full suite: CI runs
it before merge. A write refers only to landed nodes and its own change set's; pass
`--change <id>` for another change set.

## Who records a pick

The run's capture records the owner's AskUserQuestion picks after you exit, so the planner
writes no Decision for a plain owner pick. Write one only when the capture's wording would be
wrong: a proxy such as an orchestrator or arbiter answered, or the answer carries a ruling
the capture cannot word. Give it the proxy's wording and copy the asked question
verbatim into its `question` field (`"wording=<option>, picked by the orchestrator"`,
`"question=<as asked>"`): the capture skips a pick whose question a Decision already records.

## Decisions

A Decision from a typed answer quotes the owner's typed words in `quote` and links
`derived-from` the interview's Source. Decisions captured from picks
(`rhiza session interview`, after you exit) link it on their own when the Source exists in
the session's change set by then; don't write those again.
