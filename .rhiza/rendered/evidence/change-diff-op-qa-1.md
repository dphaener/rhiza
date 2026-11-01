# change-diff-op-qa-1

evidence · `01M6C6AJY9E4W9ZJE7W0NJHFDH` · actor [claude-code-8f5c32fe](../actor/claude-code-8f5c32fe.md)

## result

PASS

## statement

At 500c4d63, built go build -o $TMPDIR/rhiza-qa ./cmd/rhiza and ran it by hand in a fresh scratch repo with RHIZA_* unset. In one change set: set thing statement a->b, unset acceptance, set statement again ->c, retire gone --reason obsolete, add fresh under thing, link fresh relates-to thing. change diff printed 'updated capability thing / acceptance: ok -> unset / statement: a -> c' (one before/after for two sets), 'retired capability gone: obsolete', 'created capability fresh' with every field and no old value, and 'added link relates-to fresh thing'. After landing acceptance on thing in another change set (--record), the diff still showed only this change set's 'statement: a -> c' and 'acceptance: ok -> unset'. A pure mv of a node printed 'moved capability solo' with '@key: z -> r'; unlink printed 'retired link relates-to parent2 solo'. change diff --json and api call change.diff '{}' were cmp-identical. go test -run TestScripts/(change-diff|describe-golden) passes.

## Links

### Out

- verifies [change-diff-op](../capability/rhiza/change-sets/change-diff/change-diff-op.md)

## History

- `01M6BQJQP3M6K5RAGC5VRSK7GJ` rhiza change diff and the change-show app show each node's fields before and after, against the change set's base
