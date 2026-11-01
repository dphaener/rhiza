# change-diff-text

capability · `01M6BQKJ124RYHHZ3TAR972EA8` · parent [change-diff](../change-diff.md)

## acceptance

A Go test of the line diff and a change-diff.txtar step: a three-line statement with its middle line changed prints exactly one - line and one + line around the unchanged lines

## statement

A multi-line string field in a change diff is shown as a line diff, removed lines with - and added lines with +

## Links

### Out

- decided-by [diff-against-op-base](../../../../decision/diff-against-op-base.md)

### In

- verifies [change-diff-text-verified](../../../../evidence/change-diff-text-verified.md)
- verifies [change-diff-text-qa-1](../../../../evidence/change-diff-text-qa-1.md)
- verifies [pr-179-field-level-diff](../../../../evidence/pr-179-field-level-diff.md)

## History

- `01M6BQJQP3M6K5RAGC5VRSK7GJ` rhiza change diff and the change-show app show each node's fields before and after, against the change set's base
