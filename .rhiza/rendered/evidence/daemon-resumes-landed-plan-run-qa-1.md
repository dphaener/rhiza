# daemon-resumes-landed-plan-run-qa-1

evidence · `01M6MQYPZVQQ0QHPBEA0BPP88P` · actor [claude-code-5fde6012](../actor/claude-code-5fde6012.md)

## result

PASS

## statement

QA at dfff1995 with a binary built from the branch (go build -trimpath ./cmd/rhiza), in a scratch git repo with RHIZA_*/REX_* unset and a fake rex on PATH. Started a one-step wait workflow on a record change set, rewrote its journal's workflow to plan, ran rhiza land, then made the stop an interruption (dropped the stopped entry, appended a dead pid). rhiza daemon --permission-mode manual --once printed 'Resumed run <id>', journaled one daemon-resumed entry, and its real detached rhiza run resume --foreground process (no test hook) recorded its own process and ran the run to succeeded. Interrupted again (dead pid after daemon-resumed), a second daemon --once resumed nothing ('resumed': []) and daemon-resumed stayed at 1. A failed plan run on the landed record was not resumed. A ship-change run interrupted on a proposed, accepted and landed change set that never reached land, and a plan run interrupted on a withdrawn change set, were both left unresumed by daemon --once. Scoped tests: go test ./internal/daemon/... and the daemon-resume* testscripts pass.

## Links

### Out

- verifies [daemon-resumes-landed-plan-run](../capability/rhiza/orchestration/daemon/daemon-resumes-landed-plan-run.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
