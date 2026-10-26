# Genesis

Genesis keeps a product model in the repo: intents, capabilities, constraints, decisions,
questions, evidence and sources, joined by typed links. Every change is an entry in an
append-only log under `.genesis/log/`, grouped into change sets that are proposed, accepted
and landed. Code and model move together: the model says what the code is for and why.

## Where you are

- `genesis status` shows the repo, the current change set and what needs attention.
- `GENESIS_CHANGE` in the environment names the change set you work in. `--change <id>` on any
  command overrides it; without either, the worktree's current change set applies.
- `genesis change show --json` shows that change set: intent, state, operations.

## Vocabulary

| People say | Genesis has |
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

You work on an accepted change set (`genesis ready` names the next one). The contract:

1. Read the slice first, before reading code:
   `genesis slice --change "$GENESIS_CHANGE" --json`
2. Do the work.
3. Write back only revisions within the change set's intent, in one
   `genesis apply <file> --json`: every node and link the work adds or changes. One call, not many.
4. Anything outside the intent is not this change set's: open a new one for it and link it
   `discovered-from` this one, or add a question (recipe below). Don't write it here.
   `genesis change new "<intent>" --json` makes the new one current, so then run
   `genesis link <new> discovered-from <this> --change <new> --json` and
   `genesis change use <this>` to go back (`GENESIS_CHANGE`, if set, still wins).
5. Add a `Change: <id>` trailer to every commit, with the change set's ID:

   ```
   Stream CSV exports

   Change: 01M5PWGEF85Y185JZ2ZGCNH40N
   ```

6. Finish with `genesis land --json`, then `genesis render`, committed as the last commit.

`genesis propose` is for a draft you opened yourself with `genesis change new`.

## Reading the model

- `genesis slice <ref> --json`: the composed context for a task, the first thing to read.
- `genesis show <ref> --json`: one node, its fields and links.
- `genesis why <ref> --json`: the decisions, change sets and sources behind a node.
- `genesis log <ref> --json`: every change set that touched it.
- `genesis tree --json` and `genesis ls --json`: the capability tree; filtered lists.

Read before you write: find the capability your work belongs under, and the decisions and
constraints on it, so you add to the model rather than duplicate it.

## Hard rules

- Never edit `.genesis/log/` by hand. Write only through `genesis` commands.
- Always pass `--json`. Read the JSON; don't parse the markdown.
- Propose; run `genesis accept` only where policy allows you to.
- Never paraphrase an owner's decision: quote their exact words in the decision's fields.
  Only what the owner typed is theirs. An option an agent wrote (an AskUserQuestion choice, a
  "(Recommended)" pick) is never the owner's words: record it as the agent's wording, accepted
  by the owner.
- When something is unknown, add a question rather than guess. Park it, don't block on it.
- Use `--dry-run` to see the operations a write would make before making it.

## Recipes

Write the file, then `genesis apply <file> --json`. Each item is one write request plus `"op"`;
`genesis describe apply --json` gives the schema. Later items may refer to handles added by
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
`genesis apply <file> --record --json` opens a record change set, applies the file and lands it
in one step. Use it for capabilities the code already has, not for new intentions.

Single writes exist too: `genesis add`, `genesis set`, `genesis mv`, `genesis retire`,
`genesis link`, `genesis unlink`. Prefer one `apply` when there is more than one.

## Change sets

- `genesis change new "<intent>" --json` opens an intention change set and makes it current.
- `genesis propose`, `genesis accept`, `genesis reject`, `genesis withdraw`, `genesis land`
  move it through its lifecycle.
- `genesis change ls --json` lists them all.

## Errors

Exit codes: 0 ok, 1 error, 2 bad usage, 3 conflict, 4 not found. With `--json` an error is
`{"error":{"code","message","next"}}`. Read `next` and run it; it is the command that fixes
or explains the problem. `genesis describe <command> --json` gives any command's arguments and output.

## More

`genesis guide <topic>` prints a topic; `genesis guide ls` lists them.
Topics: `change-sets` (lifecycle, intention vs record, `--record`), `carving` (moving nodes
between change sets), `apply` (the file format), `sources` (recording a conversation). To turn
a raw conversation into a Source, use the genesis-scrub skill (`genesis guide scrub`).
