# Plan review

A fresh review of a plan written into the model. You are given only a change set ID and one
role: reviewer, critic or reviser. Each role ends with `rhiza session done`, run as its own
command only once the recording command (`rhiza run findings`, `rhiza run verdict` or the
reviser's `rhiza apply`) has succeeded. Never chain it after one with `;`: a rejected command
would still end the session with nothing recorded.

## Find the plan

The plan's change sets are that one and every change set serving an intent it adds.
`rhiza change show <id> --json` lists its ops: each intent it creates, then
`rhiza show <intent> --change <id> --json` names the change sets that serve it
(`links_in` of type `serves`). Read each with `rhiza slice --change <cs> --json` and
`rhiza why <node> --change <cs> --json`.

## Reviewer

Record what is wrong, one finding per problem, each of one of five kinds:

- `missing-acceptance`: a leaf with no `acceptance` check, or one nothing could run, or one that runs the repo's full suite
  rather than the tests covering that leaf.
- `missing-provenance`: a Decision with no `quote`, `wording` or `derived-from` Source.
- `unreal-depends-on`: a `depends-on` link where neither needs the other's work.
- `outside-intent`: a node outside its change set's intent.
- `contradiction`: nodes that contradict each other, the model or the code.

For each plan change set run `rhiza overlaps --change <cs> --json`: it lists the in-flight
change sets that write the same node or link. Record each overlap that contradicts the plan as
a `contradiction` finding naming the other change set; skip one that fits the plan.

Record the round from stdin with a heredoc, no temp file:

```
rhiza run findings - --json <<'EOF'
{"findings":[{"node":"<handle or ID>","kind":"missing-acceptance","reason":"<why>"}]}
EOF
```

The command records the round; findings are numbered f1, f2, ...
When the plan is clean, record `{"findings":[]}`: every review records a round.

## Critic

Read the latest round in `rhiza run status --json`, check each finding against the model
and the code, and confirm or drop it. Never add a finding.

```
rhiza run verdict - --json <<'EOF'
{"verdicts":[{"finding":"f1","verdict":"confirmed","reason":"<why>"}]}
EOF
```

`verdict` is `confirmed` or `dropped`.

## Reviser

Fix each confirmed finding of the latest round in the model, nothing else: one
`rhiza apply <file> --change <cs> --json` per change set holding the nodes. A node in a
landed record is fixed with a record write: `rhiza apply <file> --record --json`.
Never invent an owner's words: a missing quote becomes a question, not a made-up quote.
