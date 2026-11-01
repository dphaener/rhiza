# daemon-unjournaled-retry-qa-1

evidence · `01M6AJQ9NDDPRBNSHQC3SGANFV` · actor [claude-code-fb31714e](../actor/claude-code-fb31714e.md)

## result

PASS

## statement

Built 78aa79c0 to a scratch binary and ran the real `rhiza daemon --permission-mode manual` (no --once) in a scratch repo with two accepted milestone change sets a and b, with a fake rex on PATH that only logs (so no run journals). Both started at 03:17:08Z; after the one-minute grace `rhiza daemon status` showed each as 'Not journaled: ... started at 2026-11-01T03:17:08Z wrote no journal; retried at 2026-11-01T03:32:08Z, or sooner when a run of it journals'. Withdrawing b at 03:19:51Z dropped it from the unjournaled list at the next poll (03:19:55Z) with no new start of it. With the daemon never restarted, rex.log showed no new start until 03:32:10Z, when the daemon started a again (the first poll past 15 minutes); a went back to 'started, its run not journaled yet', and after the grace was held anew with 'retried at 2026-11-01T03:47:10Z'. b was never started again.

## Links

### Out

- verifies [daemon-unjournaled-retry](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-unjournaled-retry.md)

## History

- `01M6A0JVE5ZV75XEMB0ZNYN8SR` A daemon start that wrote no journal is retried after a hold, without restarting the daemon
