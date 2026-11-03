# run-log-verified-2

evidence · `01M6GMV31KZNA8HBQ41TSEC7DZ` · actor [claude-code-272dc04c](../actor/claude-code-272dc04c.md)

## result

PASS

## statement

Verifier at e12d4afc (attempt 2, the review fix commit). I read the diff: a daemon start now detaches into <cs>.start-<n>.log named to the child in RHIZA_RUN_LOG, and run start --foreground renames it to <run>.log before it journals. A request's plan run logs to <run>.log under q.Run. run log reads <run>.log only. This addresses the confirmed qa f1 (shared per-change-set log) and quality f2 (plan-<request>.log never found). I ran go test ./cmd/rhiza -count=1 -run TestScripts/(run-log|daemon-headless-runs|plan-requests|daemon-repo-session|headless-run-agent-tabs|run-start-queue|run-continue-queue): ok. In run-log.txtar, rhiza run log pages with --limit/--cursor/--all and exits 2 with no log. A second daemon start of the same change set logs only its own lines ('step build/build: entered'), and the first run keeps only its own. plan-requests shows run log <plan run> --all exits 0 with no plan-*.log left. go test ./internal/tui/liveapp -run Log: TestLKeyOpensSelectedRunLogLayer PASS (fake runner, l key opens a layer running rhiza run log for the selected run). In make check, fmt, vet and all go tests pass. Only self-check fails, on a stale .rhiza/rendered (actor claude-code-4d7b2408, index.md). The stale render comes from the earlier evidence commit 56a61765, not from this leaf's diff.

## Links

### Out

- verifies [run-log](../capability/rhiza/orchestration/daemon/run-log.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
