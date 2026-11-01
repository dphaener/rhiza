# change-diff-text-qa-1

evidence · `01M6C6ANKYZPB2VDGRZQMKM3GT` · actor [claude-code-8f5c32fe](../actor/claude-code-8f5c32fe.md)

## result

PASS

## statement

At 500c4d63 with a scratch build in a fresh repo: a landed four-line statement l1/l2/l3/l4 set to l1/L2/l3/l4/l5 in a change set gave, under '  statement:', exactly '      l1', '    - l2', '    + L2', '      l3', '      l4', '    + l5', so one - and one + for the changed middle line around the unchanged lines; single-line fields kept the 'a -> c' form. go test ./internal/changeset passes.

## Links

### Out

- verifies [change-diff-text](../capability/rhiza/change-sets/change-diff/change-diff-text.md)

## History

- `01M6BQJQP3M6K5RAGC5VRSK7GJ` rhiza change diff and the change-show app show each node's fields before and after, against the change set's base
