# change-show-diff-verified

evidence · `01M6C667373BJ2DGT4961KAKSY` · actor [claude-code-f9b17c44](../actor/claude-code-f9b17c44.md)

## result

PASS

## statement

At d0f1d260: read the diff (ChangeShowResponse gains changes, built by buildChangeDiff over changeset.BuildDiff, and Human ends with a '## Changes' section via ChangeDiffResponse.Human). change-diff.txtar sets thing's statement a -> c and asserts change show prints '## Changes' and '^  statement: a -> c$'; tui-write-apps.txtar cmp's the change-show app's --once output with change show --json and asserts it holds "changes" with the draft handle. Ran go test ./... (all ok), go test ./cmd/rhiza -run 'TestScripts/(change-diff|tui-write-apps|run-review-check-landed|describe-golden)$' -count=1 (ok), and make self-check fmt-check with RHIZA_ROOT/CHANGE/RUN unset (ok).

## Links

### Out

- verifies [change-show-diff](../capability/rhiza/change-sets/change-diff/change-show-diff.md)

## History

- `01M6BQJQP3M6K5RAGC5VRSK7GJ` rhiza change diff and the change-show app show each node's fields before and after, against the change set's base
