# daemon-state-at-start-qa-1

evidence · `01M6CXVQN2S6RBZNK4HVAQW014` · actor [claude-code-4d1d9a27](../actor/claude-code-4d1d9a27.md)

## result

PASS

## statement

At 423ba9c5, by hand with a scratch build (go build -o outside PATH) in a scratch git repo, RHIZA_*/REX_* unset: accepted a change set, started the real daemon (rhiza daemon --max-runs 3 --min-free 1000T --permission-mode plan; the disk guard kept it from starting any run). daemon status showed running, the skip 'Skip me: serves no milestone' and disk low. Rebuilt the binary in place; the daemon logged 'binary changed, re-executing' and the same pid took the lock again. A watcher saved every version of .git/rhiza/daemon.state.json: the new image's first write had the new start 01:26:16.082 with the old image's last poll 01:26:15.741, the skip, max runs 3 and disk_low, and a zero stopped time; its own first poll then wrote last poll 01:26:16.099. daemon status after the re-exec showed running, Started 01:26:16, the skip, no stopped line; after SIGINT it showed 'Daemon not running' with 'Last daemon ... stopped'. go test ./internal/api -run TestDaemon -count=1 passes.

## Links

### Out

- verifies [daemon-state-at-start](../capability/rhiza/orchestration/daemon/daemon-reexec-new-binary/daemon-state-at-start.md)

## History

- `01M6CA9QMHN6N95AS3F935ZVN2` daemon status keeps the daemon's state through a re-exec, and a terminal current change set comes with the command that clears it
