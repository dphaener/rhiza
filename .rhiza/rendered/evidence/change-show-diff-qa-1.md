# change-show-diff-qa-1

evidence · `01M6C6AR4VEBEG37M4HGARQV76` · actor [claude-code-8f5c32fe](../actor/claude-code-8f5c32fe.md)

## result

PASS

## statement

At 500c4d63 with a scratch build in the same scratch repo: change show ended with a '## Changes' section holding the same text as change diff, including 'statement: a -> c' for the updated statement and the line diff. rhiza tui app change-show --once was cmp-identical to change show --json and held the changes. go test -run TestScripts/tui-write-apps passes.

## Links

### Out

- verifies [change-show-diff](../capability/rhiza/change-sets/change-diff/change-show-diff.md)

## History

- `01M6BQJQP3M6K5RAGC5VRSK7GJ` rhiza change diff and the change-show app show each node's fields before and after, against the change set's base
