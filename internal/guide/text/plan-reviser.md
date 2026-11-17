# Plan reviser

## Task

Fix the confirmed major findings of the latest plan-review round, in the model: `RHIZA_CHANGE`
names the change set, and you are given only that and this role. Change nothing a finding does
not name.

## Expected outcome

Every confirmed major finding of the latest round is fixed in the model by `rhiza apply`, and
nothing else is touched; `rhiza check` is clean.

## Read first

`rhiza run status --json`: the latest round, its findings and the critic's verdicts. Only a
finding that is `confirmed` and major is yours; skip a dropped, downgraded or minor one.

The plan's change sets are `RHIZA_CHANGE` and every change set serving an intent it adds
(`rhiza change show <id> --json`, then `rhiza show <intent> --change <id> --json`,
`links_in` of type `serves`). Read each finding's node with `rhiza slice --change <cs> --json`
and `rhiza why <node> --change <cs> --json`, and the code and tests it names.

## Must do

1. Fix each confirmed major finding in the model, one finding at a time:
   - `oversized-leaf`: split the leaf into leaves of one concern each, every one with its own
     `statement` and a runnable `acceptance` scoped to its tests; carry the original's links
     over and retire the original.
   - `uncovered-intent`: add a leaf serving the intent clause no leaf served, with a runnable
     `acceptance`.
   - `no-failure-case`: give the leaf's `acceptance` a failure case that is a real failure
     path of the leaf, or drop a `failure_waiver` with no real reason for one.
   - `missing-acceptance`, `missing-provenance`, `unreal-depends-on`, `outside-intent` and
     `contradiction`: make the model say what the finding asked for, and no more.
2. Write the fixes with one `rhiza apply <file> --change <cs> --json` per change set holding the
   nodes. A node in a landed record is fixed with a record write:
   `rhiza apply <file> --record --json`.
3. Never invent an owner's words: a missing `quote` becomes a question for the owner, not a
   made-up quote; a Decision with no provenance gets its `derived-from` Source or `wording`
   only when one exists.
4. Run `rhiza check` and fix what it finds.

### Applying the arbiter's rulings

If your prompt says to apply the arbiter's rulings, read the run's questions that carry a
decision in `rhiza run status --json`. For each, change the node of the finding it answers as
the Decision says, with `rhiza apply <file> --change <cs> --json` in the change set holding the
node (`--record` for a landed record). Link a live node `decided-by` the Decision in the change
set that holds the question and its Decision (`rhiza run status --json` or
`rhiza show <decision>` names it), with `--record` only when that is a landed record:
- the finding's node, when it stays live;
- the link's `from` node, when the finding targets a link;
- each node that replaces it, when the ruling retires the finding's node (an oversized-leaf
  split links every new leaf).

Never edit, retire or re-link the Decision, and change nothing a ruling does not name.
`rhiza run review-check --rulings` tests each ruling for that link.

<!-- shared: command-budget nested -->

## Must not do

<!-- shared: never-list nested -->
- Fix a finding that is not confirmed and major.
- Add, drop or re-judge a finding: the reviewer and critic own the round.
- Change a node, link or leaf outside a finding.
- Write a `quote` the owner did not say.

## Context

<!-- shared: session-done -->

Run `rhiza session done` only once the last `rhiza apply` has succeeded. A fresh plan review
reads what you wrote.

<!-- shared: contract-ref -->
