# inbox-overlap-items

capability · `01M69DHJGRNQNZK1JFZY4QYWZ3` · parent [inbox](../inbox.md)

## acceptance

A case in inbox-list.txtar shows an overlap item in text and --json, and inbox open on it exits 2; the daemon's inbox notice labels it overlap; the inbox app groups overlap items, with a test; describe-golden.txtar is regenerated

## statement

rhiza inbox and the inbox app list each overlap between in-flight change sets as an overlap item with no run, naming both change sets and the target

## Links

### Out

- decided-by [overlaps-reuse-model](../../../../decision/overlaps-reuse-model.md)

### In

- verifies [inbox-overlap-items-verified](../../../../evidence/inbox-overlap-items-verified.md)
- verifies [inbox-overlap-items-qa-1](../../../../evidence/inbox-overlap-items-qa-1.md)
- verifies [pr-149-genesis-detects-conflicting](../../../../evidence/pr-149-genesis-detects-conflicting.md)

## History

- `01M5SWTNWB9GEJ5AA8R7K867KB` Genesis detects conflicting decisions across change sets in flight and surfaces them to the owner.
