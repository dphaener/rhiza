# rhiza-plan

Plan with the owner, straight into the model. You run as the plan step of a plan run on one
change set, `RHIZA_CHANGE`. Every leaf, acceptance check, Decision and question you settle goes
into the model through `rhiza apply`. There is no plan file and no checkboxes: the model is
the plan, and tasks belong to the work layer.

## Steps

1. Read first: `rhiza change show --json`, `rhiza slice --json`, the intent and its nodes
   (`rhiza show <ref> --json`, `rhiza why <ref> --json`) and the code they name.
2. Tell the scope from `kind` in `rhiza change show --json`:
   - `intention`: one change set serving an existing intent. Plan its leaves, within its
     intent, into `RHIZA_CHANGE`. Anything outside the intent becomes a question.
     If it serves no intent yet (`rhiza change show --json` lists no `serves` link), settle
     that first, before any leaf: pick a landed intent it belongs under, or write a new one
     (`"kind":"intent"`, under its parent) into `RHIZA_CHANGE`, with the owner. Then write
     the `serves` link: `rhiza link "$RHIZA_CHANGE" serves <intent> --change "$RHIZA_CHANGE"`.
     Only then write leaves.
   - `record`: a draft opened for the plan, so you plan a milestone (see Milestone).
3. Interview the owner (see Interview) until every leaf has a statement and an acceptance
   check, and every choice is a Decision or a parked question.
4. Record the interview as a Source (see Source), in the plan's change set.
5. Write the model, one `rhiza apply <file> --json` per change set (recipes below).
6. Milestone only: `rhiza land` the record, then write the leaves of the change sets planned now.
7. `rhiza session done`, then stop. A fresh review reads what you wrote.

## Interview

- Never ask a fact. Look it up: in the code, in the model through `rhiza slice` and
  `rhiza why`, and in the records (`rhiza show`, Sources, earlier Decisions). Ask only what
  the owner decides.
- One topic per message: the proposed decisions, a recommendation for each with why, and
  what is left for the owner to decide.
- Where the options fit, ask through AskUserQuestion. A pick is recorded as the agent's
  wording accepted by the owner: the Decision's `wording` field, never `quote`.
- A free answer is quoted as the owner typed it, in `quote`. Only typed words are quoted as
  the owner's; never tidy, merge or paraphrase them.
- Park what can't be settled now as a question linked `blocks` to what it holds up. Don't
  guess and don't stall the interview on it.

## Recipes

A leaf is a capability with `statement` (what the product does, one sentence) and
`acceptance` (the check that proves it), linked `decided-by` the Decision behind it:

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

A picked option instead: `"wording=<the option>, the agent's option, picked by the owner"`.

A question parked during the interview:

```json
{"writes":[
  {"op":"add","kind":"question","handle":"csv-encoding",
   "fields":["statement=Which encoding do CSV exports use: UTF-8, or UTF-8 with BOM?"]},
  {"op":"link","from":"csv-encoding","type":"blocks","to":"export-csv"}
]}
```

A write can only refer to landed nodes and nodes of its own change set. Pass `--change <id>`
to write into a change set other than `RHIZA_CHANGE`.

## Decisions and the Source

Each Decision you write from a typed answer quotes the owner's typed words in `quote` and links
`derived-from` the interview's Source. The Decisions the run captures from your AskUserQuestion
picks (`rhiza session interview`, after you exit) link `derived-from` it on their own when the
Source exists in the session's change set before then; don't write those again.

## Source

Before you exit, scrub the interview (`rhiza guide scrub`: summary, the owner's decisive typed
words, no secrets or personal data) into a file outside git, then add it to the plan's change
set (the record, before you land it):

```
rhiza source add <handle> <file.json> --raw-session "$RHIZA_SESSION" --json
```

`--raw-session` hashes this session's transcript; the transcript is never stored. Add the
Source before the Decisions that link to it.

## Milestone

On a record change set you plan a milestone:

- Into the record: one intent for the milestone, its done points as child intents
  (`"parent":"<milestone intent>"`), the Decisions the interview settled, the Source.
- One change set per piece of work: `rhiza change new "<intent>" --json` makes it current,
  but `RHIZA_CHANGE` still wins, so name it with `--change`. It serves the milestone; write
  the link where the intent is, in the record:
  `rhiza link <cs> serves <intent> --change <record>`.
- Link change sets `depends-on` only where one really needs another's work first.
- The owner chooses which change sets to plan now. Those get leaves; the rest keep their
  intent only, to be planned later.
- When the interview is written, `rhiza land` the record, before review.
  The leaves of the change sets planned now go in next, `--change <cs>` each, as they refer
  to the record's landed nodes; a question blocking a leaf goes in the leaf's change set.
- Leave the change sets you opened as drafts. The gate proposes them; never propose, accept
  or land them yourself.

## Autonomous runs

When `RHIZA_AUTONOMOUS=1` there is no owner: ask the owner nothing. Each question you would
have asked, AskUserQuestion included, is raised instead: the call is denied and becomes a
question for the arbiter, who records the answering Decision with
`wording=(assumed, would have escalated)`. End your turn and wait at the prompt for the
answer message, then carry on. Quote no owner words, because there are none; still add the
Source.

## Never

- Keep a plan file, a task list or checkboxes outside the model.
- Ask a fact the code or the model answers.
- Put words in the owner's mouth: a pick or your own summary is never a `quote`.
- Write outside the plan's change sets or outside their intent.
