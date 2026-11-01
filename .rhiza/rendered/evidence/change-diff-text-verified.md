# change-diff-text-verified

evidence · `01M6C53YPP6W161M2DT29CFVR5` · actor [claude-code-ff5b0627](../actor/claude-code-ff5b0627.md)

## result

PASS

## statement

At d5dedb22: go test ./internal/changeset -run TestLineDiff passes (middle changed, same, line added, line removed, single to multi cases). go test ./cmd/rhiza -run TestScripts/change-diff passes, including the new step that sets a three-line statement one/two/three to one/TWO/three and checks the '  statement:' header, '      one', '    - two', '    + TWO', '      three' lines. Because the txtar regexes only check that each line is present, I also reproduced it by hand with a scratch build in a fresh repo: change diff printed exactly one - line (two), one + line (TWO) and the two unchanged lines (one, three) around them, nothing else for the field. Read the diff: LineDiff in internal/changeset/linediff.go is an LCS that puts removed runs before added runs, and ChangeDiffResponse.Human uses it only when the old or new value contains a newline, so single-line fields keep the 'a -> b' form. gofmt -l is clean, go vet ./... is clean and go test ./... -count=1 passes.

## Links

### Out

- verifies [change-diff-text](../capability/rhiza/change-sets/change-diff/change-diff-text.md)

## History

- `01M6BQJQP3M6K5RAGC5VRSK7GJ` rhiza change diff and the change-show app show each node's fields before and after, against the change set's base
