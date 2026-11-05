# change-show-diff

capability · `01M6BQKKV7CTFV9W0T453F3RSY` · parent [change-diff](../change-diff.md)

## acceptance

change show on a change set that updates a statement prints a Changes section with the old and new statement

## statement

rhiza change show ends with a Changes section holding the change set's field-level diff

## Links

### Out

- decided-by [diff-against-op-base](../../../../decision/diff-against-op-base.md)

### In

- verifies [change-show-diff-verified](../../../../evidence/change-show-diff-verified.md)
- verifies [change-show-diff-qa-1](../../../../evidence/change-show-diff-qa-1.md)
- verifies [pr-179-field-level-diff](../../../../evidence/pr-179-field-level-diff.md)
- verifies [pr-301-delete-the-old-rex-apps](../../../../evidence/pr-301-delete-the-old-rex-apps.md)

## History

- `01M6BQJQP3M6K5RAGC5VRSK7GJ` rhiza change diff and the change-show app show each node's fields before and after, against the change set's base
- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
