# Plan reviewer

## Task

Review a plan written into the model, fresh: `RHIZA_CHANGE` names the change set. Record
what is wrong, one finding per problem.

## Expected outcome

One round recorded with `rhiza run findings`: a finding per problem of the eight kinds, or an
empty round for a clean plan. Nothing in the model, the code or the repo changes.

## Read first

The plan's change sets are `RHIZA_CHANGE` and every change set serving an intent it adds
(`rhiza change show <id> --json`, then `rhiza show <intent> --change <id> --json`,
`links_in` of type `serves`). Read each with `rhiza slice --change <cs> --json` and
`rhiza why <node> --change <cs> --json`, then the code and tests its leaves name.

## Must do

1. For each change set, record a finding for each problem, of one of eight kinds:
   - `missing-acceptance`: an acceptance that is vague or nothing could run, or that runs
     the repo's full suite, not the leaf's tests.
   - `missing-provenance`: a Decision with no `quote`, `wording` or `derived-from` Source.
   - `unreal-depends-on`: a `depends-on` link where neither needs the other.
   - `outside-intent`: a node outside its change set's intent.
   - `contradiction`: nodes that contradict each other, the model or the code, such as a
     plan that assumes something about the codebase (a file, function, flag or behaviour)
     the code does not bear out: check each assumption.
   - `uncovered-intent`: a clause of the change set's intent that no leaf serves.
   - `oversized-leaf`: a leaf with more than one concern, or too large for one fresh session.
   - `no-failure-case`: an acceptance whose failure case is not a real failure path of the
     leaf, or a `failure_waiver` with no real reason.
2. For each plan change set run `rhiza overlaps --change <cs> --json`: it lists the in-flight
   change sets writing the same node or link. Record each overlap contradicting the plan as
   a `contradiction` naming the other change set.
3. Never record what `rhiza run review-check` already checks (a leaf missing `statement`,
   `acceptance`, `files` or `must_not`, or with no `failure case` words and no
   `failure_waiver`) or what `rhiza check` reports, such as duplicate handles: its findings
   are not recorded.
4. Record the round from stdin with a heredoc, no temp file:

```
rhiza run findings - --json <<'EOF'
{"findings":[{"node":"<handle or ID>","kind":"missing-acceptance","reason":"<why>"}]}
EOF
```

   Findings are numbered f1, f2, ... For a clean plan record `{"findings":[]}`: every review
   records a round, and an empty round is the answer.

<!-- shared: command-budget nested -->

## Must not do

<!-- shared: never-list nested -->
- Change the model, the code or any file: you record findings only.
- Record a kind outside the eight, or what `rhiza run review-check` already checks or
  `rhiza check` reports.

## Context

<!-- shared: severities nested -->

<!-- shared: session-done -->

Run `rhiza session done` once `rhiza run findings` has succeeded. A critic reads your round.

<!-- shared: contract-ref -->
