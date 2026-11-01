# milestones-in-number-order-qa-1

evidence · `01M6BAR8CMHQ4SBCZJRFW9T8YQ` · actor [claude-code-9b21f3fe](../actor/claude-code-9b21f3fe.md)

## result

PASS

## statement

QA at 0bbdfac5 with a scratch build (go build ./cmd/rhiza) in a fresh git repo, fake rex on PATH. Intents recorded in the order ten(milestone=10), plain(milestone=alpha), nine(9), two(2), zeta(5), alpha(5), each served by an accepted change set with a leaf; two more change sets serve nine, Nine B depends-on Nine A, B created first. rhiza roadmap groups: two, zeta, alpha, nine, ten, plain (numbered by number, 9 before 10, the tie at 5 by ID since neither has a key, unnumbered last); within nine, Nine A comes before Nine B. rhiza daemon --once --max-runs 2 started two then zeta; in a four-intent repo with --max-runs 1 it started nine's change set over ten's. On this repo's own model the new and installed builds give the same roadmap groups (milestone-12, milestone-13, Backlog, areas). go test ./cmd/rhiza -run TestScripts/^daemon-picks$ passes.

## Links

### Out

- verifies [milestones-in-number-order](../capability/rhiza/orchestration/daemon/daemon-picks/milestones-in-number-order.md)

## History

- `01M6ATAVKKGF9SWBWVTF7N98QS` The roadmap and the daemon order milestones by their number
