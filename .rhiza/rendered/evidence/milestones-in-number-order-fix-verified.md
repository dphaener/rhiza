# milestones-in-number-order-fix-verified

evidence · `01M6BBV7GN5C07CTHT7BB3XRGB` · actor [claude-code-f21f6cc1](../actor/claude-code-f21f6cc1.md)

## result

PASS

## statement

At df62c547 (review fixer commit): go test ./cmd/rhiza -run TestScripts/daemon-picks -count=1 passes; its milestone-order section sets up intents ten (milestone=10) and nine (milestone=9) with a ready accepted change set serving each, asserts rhiza roadmap --json lists nine before ten, and asserts rhiza daemon --once --max-runs 1 starts only the milestone-9 change set. The fix commit adds milestones-by-number supersedes m2cs2-milestone-order and milestone-groups decided-by milestones-by-number (rhiza show m2cs2-milestone-order --change shows 'superseded by milestones-by-number'; rhiza why milestone-groups lists milestones-by-number), resolving model finding f1. make check exits 0 (fmt-check, vet, go test ./..., self-check incl. render).

## Links

### Out

- verifies [milestones-in-number-order](../capability/rhiza/orchestration/daemon/daemon-picks/milestones-in-number-order.md)

## History

- `01M6ATAVKKGF9SWBWVTF7N98QS` The roadmap and the daemon order milestones by their number
