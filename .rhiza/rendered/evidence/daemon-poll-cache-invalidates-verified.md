# daemon-poll-cache-invalidates-verified

evidence · `01M6HEDFN8NFMXKTCABS193H64` · actor [claude-code-866d24a5](../actor/claude-code-866d24a5.md)

## result

PASS

## statement

At 8f349a0e I ran go test ./internal/api -run TestDaemonPollCacheInvalidates -count=1: PASS. The test polls once, then for each of an accepted change set, a daemon setting change (rhiza.autonomy.ship owner->auto), a daemonStart, a moved ref and a daemonResume, it asserts through the ctx.readLog counter that the next poll reads the log again. It also asserts the poll acts on the change: the newly accepted change set is skipped under owner, then picked once the setting is auto. The fix for attempt-1 finding f1 holds: the setting step now asserts on pickReads, the reads taken before the start's detach. In scratch worktrees I made three mutations and each failed the test: (1) leaving the rhiza.* git config output out of daemonStamp fails at daemon_fold_test.go:311 (0 reads before the pick); (2) leaving the for-each-ref output out of daemonStamp fails at :327 (ref moved); (3) removing the daemonFold invalidate in the detach path fails at :320 (start) and :336 (resume). make fmt-check and vet pass. self-check fails only on unrendered actor files, which are expected before landing. In the full go test ./..., the only failures were timeouts and timing budgets under VM load around 12 on 4 cores with 35% steal. arbiter-agent-tabs, plan-workflow and index TestBudget pass when rerun alone. api TestReadBudget fails at base ebee728a too under the same load, so it is environmental and not caused by this leaf.

## Links

### Out

- verifies [daemon-poll-cache-invalidates](../capability/rhiza/orchestration/daemon/daemon-poll-cache-invalidates.md)

## History

- `01M6H6SGKR1N5GBG7XPECKN9Z3` Fix The daemon pins a full CPU core: on the 4-core VM rhiza daemon uses ~100% CPU continuously (120 CPU-minutes in 2 hours) because each poll (daemonPoll = 5s, internal/api/daemon.go) takes 10-15 s of work: it re-reads and re-folds the whole log and recomputes picks every time. A poll should be cheap when nothing changed: skip the fold/pick work when the log, git refs, run records and settings are unchanged since the last poll (cache keyed on their stamps), and measure a poll's cost in a test. Agents' rhiza apply calls on the same VM also take minutes under this load
