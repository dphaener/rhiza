# daemon-prunes-orphan-leftovers-qa-1

evidence · `01M774R4RXVV6W08KWQWAFQ3TY` · actor [claude-code-98a70447](../actor/claude-code-98a70447.md)

## result

PASS

## statement

QA at c90d32c99 with the built binary in the same scratch root, after the run-end scenario. I removed the ended run's worktree by hand with git worktree remove, so its index DB became an orphan. I recreated that ended session's scratch folder and aged it 2 days (touch -d), and recreated the abandoned run's session folder fresh. I also added an aged folder for a session ID nobody recorded. `rhiza daemon --once --max-runs 0 --permission-mode manual` logged 'rhiza: daemon: leftovers sweep freed 120.9 KiB: 1 index databases, 1 session scratch folders'. The orphan index and the aged ended-session folder were gone. The fresh folder, the unrecorded-session folder and other/keep stayed. A new live run (stopped at its wait) whose scratch folder I aged 3 days kept its folder and its index across another daemon --once, which logged no sweep line. I did not exercise the hourly spacing by hand; it is covered by TestDaemonSweepLeftovers with an injected clock.

## Links

### Out

- verifies [daemon-prunes-orphan-leftovers](../capability/rhiza/orchestration/daemon/daemon-prunes-orphan-leftovers.md)

## History

- `01M763YK58RWNJ35P9S8BSFKVQ` Runs clean up Rhiza's own leftovers: when a run ends, Rhiza removes the scratch folders of the agent sessions it started and its worktree's index database; the daemon prunes orphans (index files whose worktree is gone, old session scratch folders) so disk use no longer grows with runs.
