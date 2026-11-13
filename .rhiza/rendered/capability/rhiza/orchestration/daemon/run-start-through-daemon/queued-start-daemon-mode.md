# queued-start-daemon-mode

capability · `01M7982F0JXM91XG2GD8WQZ75C` · parent [run-start-through-daemon](../run-start-through-daemon.md)

## acceptance

A testscript (run-start-queue.txtar or a new one) with rhiza.daemon.permissionMode set: rhiza run start <cs> --workflow plan --autonomous with no --permission-mode queues the run and the queued start records no mode; rhiza daemon --once then starts it with --permission-mode set to the daemon's setting; with the setting naming a mode the agent rejects, the queued start is refused at queue time and nothing is queued; go test -trimpath ./internal/api/... covers checkAgentSteps getting the daemon's mode

## statement

A run start queued for the daemon that names no --permission-mode is checked at queue time with rhiza.daemon.permissionMode standing in for the mode, and is queued with no mode recorded, so the daemon starts it with its own rhiza.daemon.permissionMode as it does its own picks; anything that re-checks a queued run's agent steps before the daemon starts it uses the same stand-in

## Links

### Out

- decided-by [queued-run-daemon-mode](../../../../../decision/queued-run-daemon-mode.md)

### In

- verifies [queued-start-daemon-mode-verified](../../../../../evidence/queued-start-daemon-mode-verified.md)
- verifies [queued-start-daemon-mode-qa-1](../../../../../evidence/queued-start-daemon-mode-qa-1.md)
- verifies [queued-start-daemon-mode-verified-2](../../../../../evidence/queued-start-daemon-mode-verified-2.md)
- verifies [pr-519-fix-a-run-queued-by-hand](../../../../../evidence/pr-519-fix-a-run-queued-by-hand.md)

## History

- `01M797TRSFCEHV3FXTYVRV7D74` Fix A run queued by hand (rhiza run start <cs> --workflow plan --autonomous) does not get the daemon's own permission mode: the CLI refuses --permission-mode bypassPermissions, and with none the builtin plan's interview step is refused ('agent step with no permission mode'), so a queued plan run must pick a weaker mode. On 2026-11-13 ~00:40 plan run 01M796JBG1 (change set 01M76XR5, queued with --permission-mode acceptEdits) stopped at its first shell command on a permission prompt (inbox 84109a8). rhiza run abandon then refused it ('running in another process'), so the orchestrator killed its driver and agent by PID, abandoned it, and requeued with --permission-mode auto (run 01M796PA86, works). Wanted: a run the daemon drives takes rhiza.daemon.permissionMode when the start names none (the daemon's own picks already do), so an owner or orchestrator can queue plan runs for accepted no-leaves change sets without choosing a mode.
