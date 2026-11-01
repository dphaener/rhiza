# change-diff-op-verified

evidence · `01M6C45FQ9C0Y13JJZMV1NAHHV` · actor [claude-code-89561799](../actor/claude-code-89561799.md)

## result

PASS

## statement

At 7726eaf6, go test ./cmd/rhiza -run 'TestScripts/(change-diff|describe-golden|frame)$' -count=1 passes. change-diff.txtar checks, in order: set after add gives 'statement: a -> b' and the unset 'acceptance: ok -> unset'; a second set gives one 'statement: a -> c' with no 'a -> b'; a retire prints 'retired capability gone: not needed'; a create prints 'created capability fresh' with '  statement: new one' (no old value) and 'added link relates-to fresh thing'; after another change set lands acceptance on thing with --record, the diff still shows 'statement: a -> c' and not 'landed first'; change diff --json and api call change.diff '{}' are cmp-identical. Re-running describe-golden with -update leaves the tree clean, so describe-golden.txtar is current (change.diff added). make fmt-check vet test all pass. Read the diff: internal/changeset/diff.go rebuilds the old node by replaying the base chain of other change sets' ops from the first op's base, and internal/api/change_diff.go registers change.diff with a thin Human view. make self-check in the checkout fails only on rendered actor files for sessions started after this commit (claude-code-4831f979, claude-code-89561799), not on this leaf.

## Links

### Out

- verifies [change-diff-op](../capability/rhiza/change-sets/change-diff/change-diff-op.md)

## History

- `01M6BQJQP3M6K5RAGC5VRSK7GJ` rhiza change diff and the change-show app show each node's fields before and after, against the change set's base
