# queued-run-daemon-mode

decision · `01M7981417W8MH0Y9P9680KJHE`

## reasoning

The change set's intent asks exactly for this: a run the daemon drives takes rhiza.daemon.permissionMode when the start names none, so an owner or orchestrator can queue plan runs without choosing a mode (the 2026-11-13 acceptEdits run stalled on a permission prompt). Checking with the daemon's setting at queue time keeps the existing guarantee in run.go that a start whose agent steps cannot run is refused before any git or write, while recording no mode lets daemonStartsQueued apply the setting current at start time rather than freezing a copy. A foreground start or an unset daemon setting has no mode to fall back on, so the agent step genuinely cannot run; keeping the refusal there is the safe, minimal choice, and naming the setting in the message points users at the new way out. Both halves are the question-asker's own recommendations and fit the code (checkAgentSteps at run.go:672/695, daemon.go:128 already checks with conf.mode).

## statement

A run start that names no --permission-mode and is queued for the daemon (not --foreground) is checked at queue time with rhiza.daemon.permissionMode standing in for the mode (checkAgentSteps gets the daemon's setting), and is queued with no mode recorded, so the daemon fills in its own rhiza.daemon.permissionMode when it starts the run, as it already does for its own picks. When the start names no mode and rhiza.daemon.permissionMode is unset, or the start is --foreground (no daemon drives it), today's refusal stays, and its message names rhiza.daemon.permissionMode as a third way out next to --permission-mode and the step's permission_mode. Anything that re-checks a queued run's agent steps before the daemon starts it uses the same stand-in.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-16742c4a-1](../question/question-16742c4a-1.md)
- derived-from [queued-run-mode-interview](../source/queued-run-mode-interview.md)

### In

- decided-by [queued-start-daemon-mode](../capability/rhiza/orchestration/daemon/run-start-through-daemon/queued-start-daemon-mode.md)
- decided-by [no-mode-refusal-names-setting](../capability/rhiza/orchestration/daemon/run-start-through-daemon/no-mode-refusal-names-setting.md)

## History

- `01M797TRSFCEHV3FXTYVRV7D74` Fix A run queued by hand (rhiza run start <cs> --workflow plan --autonomous) does not get the daemon's own permission mode: the CLI refuses --permission-mode bypassPermissions, and with none the builtin plan's interview step is refused ('agent step with no permission mode'), so a queued plan run must pick a weaker mode. On 2026-11-13 ~00:40 plan run 01M796JBG1 (change set 01M76XR5, queued with --permission-mode acceptEdits) stopped at its first shell command on a permission prompt (inbox 84109a8). rhiza run abandon then refused it ('running in another process'), so the orchestrator killed its driver and agent by PID, abandoned it, and requeued with --permission-mode auto (run 01M796PA86, works). Wanted: a run the daemon drives takes rhiza.daemon.permissionMode when the start names none (the daemon's own picks already do), so an owner or orchestrator can queue plan runs for accepted no-leaves change sets without choosing a mode.
