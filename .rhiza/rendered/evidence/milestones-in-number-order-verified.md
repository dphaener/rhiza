# milestones-in-number-order-verified

evidence · `01M6BAKK4HW04AE4HDCHN50D60` · actor [claude-code-3a31b0f7](../actor/claude-code-3a31b0f7.md)

## result

PASS

## statement

At 3b2eb9bd, daemon-picks.txtar now has a repo with intents nine (milestone=9) and ten (milestone=10, recorded first), each served by an accepted change set with a leaf. go test ./cmd/rhiza -run 'TestScripts/^daemon-picks$' passes: rhiza roadmap --json lists nine's group before ten's, and rhiza daemon --once --max-runs 1 starts exactly one run, Nine's change set (01M5PWGW...). With internal/api/plan.go reverted to dee5e4d7 the same script fails at daemon-picks.txtar:116 (roadmap order). The diff sorts siblings by milestone number (a JSON number or numeric string; numbered ones before unnumbered), then key, then ID. Depends-on topo order within a group is unchanged. make check exits 0 once RHIZA_ROOT/RHIZA_CHANGE/RHIZA_RUN are unset. With them set, self-check reads the run worktree's actor records and reports their rendered pages as missing, which is unrelated to this commit. Both attempt-1 findings are fixed: the stray actor render is gone and the scenario moved into daemon-picks.txtar.

## Links

### Out

- verifies [milestones-in-number-order](../capability/rhiza/orchestration/daemon/daemon-picks/milestones-in-number-order.md)

## History

- `01M6ATAVKKGF9SWBWVTF7N98QS` The roadmap and the daemon order milestones by their number
