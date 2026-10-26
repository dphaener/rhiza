# Rhiza

Rhiza keeps a product model in the repo: intents, capabilities, constraints, decisions,
questions, evidence and sources, joined by typed links. Every change is an entry in an
append-only log under `.rhiza/log/`, grouped into change sets that are proposed, accepted
and landed. Code and model move together: the model says what the code is for and why.

## Where you are

- `rhiza status` shows the repo, the current change set and what needs attention.
- `RHIZA_CHANGE` in the environment names the change set you work in. `--change <id>` on any
  command overrides it; without either, the worktree's current change set applies.
- `rhiza change show --json` shows that change set: intent, state, operations.

## Vocabulary

| People say | Rhiza has |
|---|---|
| feature, requirement, spec, epic | capability (in a tree, under a parent capability) |
| goal, objective, outcome | intent |
| rule, non-functional requirement, limit | constraint |
| ADR, design decision, choice | decision |
| open question, unknown, TBD | question |
| test result, benchmark, measurement | evidence |
| conversation, meeting notes, transcript | source |
| user, team, owner, role | actor |
| ticket, PR scope, piece of work | change set (its tasks belong to the work layer, not the model) |

Kinds: intent, capability, constraint, decision, question, evidence, source, actor, anchor.
Nodes have a handle (lower-case words joined by dashes) and fields (`statement` is the usual one).
Refer to a node by handle, ID, or a unique ID prefix of 4+ characters.

## The core loop

You work on an accepted change set (`rhiza ready` names the next one). The contract:

1. Read the slice first, before reading code:
   `rhiza slice --change "$RHIZA_CHANGE" --json`
2. Do the work.
3. Write back only revisions within the change set's intent, in one
   `rhiza apply <file> --json`: every node and link the work adds or changes. One call, not many.
4. Anything outside the intent is not this change set's: open a new one for it and link it
   `discovered-from` this one, or add a question (recipe below). Don't write it here.
   `rhiza change new "<intent>" --json` makes the new one current, so then run
   `rhiza link <new> discovered-from <this> --change <new> --json` and
   `rhiza change use <this>` to go back (`RHIZA_CHANGE`, if set, still wins).
5. Add a `Change: <id>` trailer to every commit, with the change set's ID:

   ```
   Stream CSV exports

   Change: 01M5PWGEF85Y185JZ2ZGCNH40N
   ```

6. Finish with `rhiza land --json`, then `rhiza render`, committed as the last commit.

`rhiza propose` is for a draft you opened yourself with `rhiza change new`.

## Reading the model

- `rhiza slice <ref> --json`: the composed context for a task, the first thing to read.
- `rhiza show <ref> --json`: one node, its fields and links.
- `rhiza why <ref> --json`: the decisions, change sets and sources behind a node.
- `rhiza log <ref> --json`: every change set that touched it.
- `rhiza tree --json` and `rhiza ls --json`: the capability tree; filtered lists.

Read before you write: find the capability your work belongs under, and the decisions and
constraints on it, so you add to the model rather than duplicate it.

## Hard rules

- Never edit `.rhiza/log/` by hand. Write only through `rhiza` commands.
- Always pass `--json`. Read the JSON; don't parse the markdown.
- Propose; run `rhiza accept` only where policy allows you to.
- Never paraphrase an owner's decision: quote their exact words in the decision's fields.
  Only what the owner typed is theirs. An option an agent wrote (an AskUserQuestion choice, a
  "(Recommended)" pick) is never the owner's words: record it as the agent's wording, accepted
  by the owner.
- When something is unknown, add a question rather than guess. Park it, don't block on it.
- Use `--dry-run` to see the operations a write would make before making it.

## Recipes

Write the file, then `rhiza apply <file> --json`. Each item is one write request plus `"op"`;
`rhiza describe apply --json` gives the schema. Later items may refer to handles added by
earlier ones.

**Add a capability** under an existing one. Write `statement` as one sentence saying what the
product does, not how the code does it:

```json
{"writes":[
  {"op":"add","kind":"capability","handle":"export-csv","parent":"reports",
   "fields":["statement=Users export any report as CSV"]}
]}
```

**Record a decision with its options**, and hang it off what it decides:

```json
{"writes":[
  {"op":"add","kind":"decision","handle":"csv-streaming",
   "fields":["statement=Stream CSV rows instead of building the file in memory",
             "quote=\"Stream it, reports get huge\" (owner, 2026-10-20)"],
   "json_fields":["options=[\"stream rows\",\"build in memory\",\"background job\"]"]},
  {"op":"link","from":"export-csv","type":"decided-by","to":"csv-streaming"}
]}
```

**Raise a question** instead of guessing:

```json
{"writes":[
  {"op":"add","kind":"question","handle":"csv-encoding",
   "fields":["statement=Which encoding do CSV exports use: UTF-8, or UTF-8 with BOM for Excel?"]},
  {"op":"link","from":"csv-encoding","type":"blocks","to":"export-csv"}
]}
```

When it is settled, link the decision to it: `{"op":"link","from":"<decision>","type":"answers","to":"csv-encoding"}`.
A decision that came from a conversation links to its source with `derived-from`.

**Record existing code** (write down what already is, no review needed):
`rhiza apply <file> --record --json` opens a record change set, applies the file and lands it
in one step. Use it for capabilities the code already has, not for new intentions.

Single writes exist too: `rhiza add`, `rhiza set`, `rhiza mv`, `rhiza retire`,
`rhiza link`, `rhiza unlink`. Prefer one `apply` when there is more than one.

## Change sets

- `rhiza change new "<intent>" --json` opens an intention change set and makes it current.
- `rhiza propose`, `rhiza accept`, `rhiza reject`, `rhiza withdraw`, `rhiza land`
  move it through its lifecycle.
- `rhiza change ls --json` lists them all.

## Errors

Exit codes: 0 ok, 1 error, 2 bad usage, 3 conflict, 4 not found. With `--json` an error is
`{"error":{"code","message","next"}}`. Read `next` and run it; it is the command that fixes
or explains the problem. `rhiza describe <command> --json` gives any command's arguments and output.

## More

`rhiza guide <topic>` prints a topic; `rhiza guide ls` lists them.
Topics: `change-sets` (lifecycle, intention vs record, `--record`), `carving` (moving nodes
between change sets), `apply` (the file format), `sources` (recording a conversation). To turn
a raw conversation into a Source, use the rhiza-scrub skill (`rhiza guide scrub`).
