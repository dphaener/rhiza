# run-worktree-carries-accept-verified

evidence · `01M6S5ZX4QY790PCGW1CVAM66V` · actor [claude-code-aa46bbcf](../actor/claude-code-aa46bbcf.md)

## result

PASS

## statement

Verified at 1df154a8. `go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/run-worktree-carries-accept` passes. The script commits a proposed change set with leaves, accepts it in the main clone (entry untracked), starts `rhiza run start --foreground` in a new worktree, then a second change set with its worktree reused on its branch, accepted between runs. In each, `change show` reports accepted, the accept entry is present and untracked, `rhiza ship land` lands with a clean tree, and the accept entry file is in the pushed branch's tree on the bare origin. With internal/api/session.go reverted to 5ca57f63, the script fails at line 39 (worktree reports proposed), so it covers the fix. Diff read: worktreeWorkspace now copies the main clone's change set entries whenever the main clone has that change set's log, never overwriting (copyNewFile), then reopens. `go test -trimpath -count=1 ./internal/api/... -skip TestReadBudget` passes. TestReadBudget failed timing budgets (different reads each run) at 1df154a8 and equally at the base 5ca57f63 under load average 6-14; it times read paths this change does not touch, so the failure was there before this change and comes from machine load.

## Links

### Out

- verifies [run-worktree-carries-accept](../capability/rhiza/orchestration/shipping/run-worktree-carries-accept.md)

## History

- `01M6RXNQFZX8YM80WFZJCXBRAJ` Fix Accepting a change set from the Inbox app never reaches the run that builds it. On 2026-11-06 the owner accepted Import procedure (01M6J36X) in the Inbox app at 14:31:10Z; the accept entry was written to the main clone's .rhiza/log but stayed uncommitted and unpushed. The daemon started the ship-change run 11 s later on a branch from origin, without the entry, and the run built, passed review and CI, then failed at land: 'cannot ship change set 01M6J36X...: it is proposed, not accepted' (the orchestrator copied the entry onto the run branch by hand to get it merged as #328). An accept from the Inbox app (and rhiza accept on the main clone generally) must reach the run's branch, e.g. by committing and pushing the accept to main or by the run carrying it, so the run never sees an accepted change set as proposed.
