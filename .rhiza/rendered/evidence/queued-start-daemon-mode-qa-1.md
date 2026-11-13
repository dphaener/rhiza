# queued-start-daemon-mode-qa-1

evidence · `01M799H774CTS52H3DQRXBCJS7` · actor [claude-code-82543a16](../actor/claude-code-82543a16.md)

## result

PASS

## statement

QA at 30e69b6fe with a binary built from the branch (go build -trimpath ./cmd/rhiza), in a scratch Rhiza root and repo with fake rex and claude on PATH and a clean env. git config rhiza.daemon.permissionMode manual, then rhiza run start --workflow plan --autonomous with no --permission-mode: 'Queued a plan run ... for the daemon', and .git/rhiza/starts.json holds the start with no permission_mode key. rhiza daemon --max-runs 3 --once (no --permission-mode): 'Started a run of <cs>', the queue drained, and the run's journal records permission_mode "manual" and autonomous true; run ls shows it live at step interview. With the setting = nonsense the start is refused at queue time (claude-code: unknown permission mode "nonsense", exit 2) and nothing is queued. With the setting = bypassPermissions the start is refused by the dedicated-machine trust check (exit 1, nothing queued) exactly like an explicit --permission-mode bypassPermissions; after writing trust = "dedicated" to machine.toml the same start queues with no mode recorded, which is the incident's case.

## Links

### Out

- verifies [queued-start-daemon-mode](../capability/rhiza/orchestration/daemon/run-start-through-daemon/queued-start-daemon-mode.md)

## History

- `01M797TRSFCEHV3FXTYVRV7D74` Fix A run queued by hand (rhiza run start <cs> --workflow plan --autonomous) does not get the daemon's own permission mode: the CLI refuses --permission-mode bypassPermissions, and with none the builtin plan's interview step is refused ('agent step with no permission mode'), so a queued plan run must pick a weaker mode. On 2026-11-13 ~00:40 plan run 01M796JBG1 (change set 01M76XR5, queued with --permission-mode acceptEdits) stopped at its first shell command on a permission prompt (inbox 84109a8). rhiza run abandon then refused it ('running in another process'), so the orchestrator killed its driver and agent by PID, abandoned it, and requeued with --permission-mode auto (run 01M796PA86, works). Wanted: a run the daemon drives takes rhiza.daemon.permissionMode when the start names none (the daemon's own picks already do), so an owner or orchestrator can queue plan runs for accepted no-leaves change sets without choosing a mode.
