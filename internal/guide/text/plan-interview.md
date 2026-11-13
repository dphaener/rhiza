# Plan interview

## Interview

- Never ask a fact. Look it up: in the code, in the model through `rhiza slice` and
  `rhiza why`, and in the records (`rhiza show`, Sources, earlier Decisions). Ask only what
  the owner decides.
- One topic per message: the proposed decisions, a recommendation for each with why, and
  what is left for the owner to decide.
- Where the options fit, ask through AskUserQuestion. The capture records picks as the
  agent's wording accepted by the owner: the Decision's `wording`, never `quote`.
- A free answer is quoted as the owner typed it, in `quote`. Only typed words are quoted as
  the owner's; never tidy, merge or paraphrase them.
- Park what can't be settled now as a question linked `blocks` to what it holds up:

```json
{"writes":[
  {"op":"add","kind":"question","handle":"csv-encoding",
   "fields":["statement=Which encoding do CSV exports use: UTF-8, or UTF-8 with BOM?"]},
  {"op":"link","from":"csv-encoding","type":"blocks","to":"export-csv"}
]}
```

## Source

Before you exit, scrub the interview (`rhiza guide scrub`: summary, the owner's decisive typed
words, no secrets or personal data) into a file outside git, then add it to the plan's change
set (the record, before you land it), before the Decisions that link to it:

```
rhiza source add <handle> <file.json> --raw-session "$RHIZA_SESSION" --json
```

`--raw-session` hashes this session's transcript; the transcript is never stored.

## Milestone

On a `record` change set you plan a milestone:

- Into the record: one intent for the milestone, its done points as child intents
  (`"parent":"<milestone intent>"`), the Decisions the interview settled, the Source.
- One change set per piece of work: `rhiza change new "<intent>" --json` makes it current,
  but `RHIZA_CHANGE` still wins, so name it with `--change`. It serves the milestone; write
  the link in the record: `rhiza link <cs> serves <intent> --change <record>`.
- Name the milestone intent from `RHIZA_MILESTONE_NUMBER` (N): handle `milestone-N`, field
  `milestone=N`, its done points' handles prefixed `mN-`. Never count milestones from your
  own branch.
- Link `depends-on` only where one needs another's work first.
- The owner chooses which change sets to plan now. Those get leaves; the rest keep their
  intent only.
- When the interview is written, `rhiza land` the record, before review. The leaves of the
  change sets planned now go in next, `--change <cs>` each; a question blocking a leaf goes
  in the leaf's change set.
- Leave the change sets you opened as drafts: the gate proposes them.
