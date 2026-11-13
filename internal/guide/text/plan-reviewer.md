# Plan reviewer

## Task

Review a plan written into the model, fresh: `RHIZA_CHANGE` names the change set. Record
what is wrong, one finding per problem, and nothing else.

## Expected outcome

One round recorded with `rhiza run findings`: a finding for each problem of the seven kinds,
or an empty round when the plan is clean. Nothing in the model, the code or the repo changes.

## Read first

The plan's change sets are `RHIZA_CHANGE` and every change set serving an intent it adds
(`rhiza change show <id> --json`, then `rhiza show <intent> --change <id> --json`,
`links_in` of type `serves`). Read each with `rhiza slice --change <cs> --json` and
`rhiza why <node> --change <cs> --json`, then the code and tests its leaves name.

## Must do

1. For each change set, check the plan against this list and record a finding for each
   problem, of one of seven kinds:
   - `missing-acceptance`: a leaf with no `acceptance`, or one that is vague or nothing could
     run, or one that runs the repo's full suite rather than the tests covering that leaf.
   - `missing-provenance`: a Decision with no `quote`, `wording` or `derived-from` Source.
   - `unreal-depends-on`: a `depends-on` link where neither needs the other's work.
   - `outside-intent`: a node outside its change set's intent.
   - `contradiction`: nodes that contradict each other, the model or the code. A plan that
     assumes something about the codebase (a file, function, flag or behaviour) that the
     code does not bear out is a contradiction: check each such assumption against the code.
   - `uncovered-intent`: a clause of the change set's intent that no leaf serves.
   - `oversized-leaf`: a leaf with more than one concern, or too large for one fresh session.
2. For each plan change set run `rhiza overlaps --change <cs> --json`: it lists the in-flight
   change sets that write the same node or link. Record each overlap that contradicts the
   plan as a `contradiction` finding naming the other change set; skip one that fits the plan.
3. Never record what `rhiza check` already reports: its findings are not recorded.
4. Record the round from stdin with a heredoc, no temp file:

```
rhiza run findings - --json <<'EOF'
{"findings":[{"node":"<handle or ID>","kind":"missing-acceptance","reason":"<why>"}]}
EOF
```

   Findings are numbered f1, f2, ... When the plan is clean, record `{"findings":[]}`: every
   review records a round, so an empty round is the answer for a clean plan.

<!-- shared: command-budget nested -->

## Must not do

<!-- shared: never-list nested -->
- Change the model, the code or any file: you record findings only.
- Record a finding of a kind outside the seven, or one `rhiza check` already reports.

## Context

<!-- shared: severities nested -->

<!-- shared: session-done -->

Run `rhiza session done` only once `rhiza run findings` has succeeded. A critic reads your
round.

<!-- shared: contract-ref -->
