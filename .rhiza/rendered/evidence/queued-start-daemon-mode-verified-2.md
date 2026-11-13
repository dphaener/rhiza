# queued-start-daemon-mode-verified-2

evidence · `01M799YC6NFH40V0Z39ATY56CT` · actor [claude-code-f7441861](../actor/claude-code-f7441861.md)

## result

PASS

## statement

Verifier at 6937db290: go test -trimpath ./cmd/rhiza -run 'TestScripts/run-start-queue$' -count=1 passed. The script sets rhiza.daemon.permissionMode=acceptEdits, queues rhiza run start --workflow plan --autonomous with no --permission-mode, asserts starts.json records no permission_mode, then rhiza daemon --once starts it and the script waits for a run journal recording "permission_mode":"acceptEdits" (a mode no earlier run used), closing review finding f1; with the setting unset or nonsense the queued start is refused (stderr names the setting / mode) and the unset case asserts nothing queued; --foreground gets no stand-in. Read the diff: queuedCheckMode (internal/api/run.go) feeds rhiza.daemon.permissionMode to both checkAgentSteps calls before queueRunStart, so a refusal writes nothing, and the queued start still records r.PermissionMode (empty); daemonQueuedStarts falls back to the daemon's mode; other checkAgentSteps callers re-check started runs whose journal holds the daemon's mode. go test -trimpath ./internal/api/... -run 'TestQueuedCheckMode|TestRunStartAgentRefusals|TestDaemon' -count=1 passed, incl. TestQueuedCheckModeStandsInDaemonMode (queued gets the daemon's mode, named mode wins, foreground none).

## Links

### Out

- verifies [queued-start-daemon-mode](../capability/rhiza/orchestration/daemon/run-start-through-daemon/queued-start-daemon-mode.md)

## History

- `01M797TRSFCEHV3FXTYVRV7D74` Fix A run queued by hand (rhiza run start <cs> --workflow plan --autonomous) does not get the daemon's own permission mode: the CLI refuses --permission-mode bypassPermissions, and with none the builtin plan's interview step is refused ('agent step with no permission mode'), so a queued plan run must pick a weaker mode. On 2026-11-13 ~00:40 plan run 01M796JBG1 (change set 01M76XR5, queued with --permission-mode acceptEdits) stopped at its first shell command on a permission prompt (inbox 84109a8). rhiza run abandon then refused it ('running in another process'), so the orchestrator killed its driver and agent by PID, abandoned it, and requeued with --permission-mode auto (run 01M796PA86, works). Wanted: a run the daemon drives takes rhiza.daemon.permissionMode when the start names none (the daemon's own picks already do), so an owner or orchestrator can queue plan runs for accepted no-leaves change sets without choosing a mode.
