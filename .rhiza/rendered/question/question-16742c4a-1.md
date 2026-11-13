# question-16742c4a-1

question · `01M797ZMWVQ0XYXRPBN7G7FE8C`

## session

16742c4a-fa0f-4ac5-b6a2-01a1d27db9cb

## statement

Change set 01M797T (queued run takes the daemon's permission mode): rhiza run start without --permission-mode, queued for the daemon (not --foreground), today fails checkAgentSteps at queue time. I recommend: at queue time, check the agent steps with rhiza.daemon.permissionMode in its place, and queue the start with no mode so the daemon fills in its setting when it starts the run (daemonStartsQueued already does m = mode). Which? / Change set 01M797T: when the start names no mode and rhiza.daemon.permissionMode is unset (or the start is --foreground, which no daemon drives), what happens? I recommend keeping today's refusal, with the message naming rhiza.daemon.permissionMode as a third way out next to --permission-mode and the step's permission_mode.

## Links

### In

- answers [queued-run-daemon-mode](../decision/queued-run-daemon-mode.md)

## History

- `01M797TRSFCEHV3FXTYVRV7D74` Fix A run queued by hand (rhiza run start <cs> --workflow plan --autonomous) does not get the daemon's own permission mode: the CLI refuses --permission-mode bypassPermissions, and with none the builtin plan's interview step is refused ('agent step with no permission mode'), so a queued plan run must pick a weaker mode. On 2026-11-13 ~00:40 plan run 01M796JBG1 (change set 01M76XR5, queued with --permission-mode acceptEdits) stopped at its first shell command on a permission prompt (inbox 84109a8). rhiza run abandon then refused it ('running in another process'), so the orchestrator killed its driver and agent by PID, abandoned it, and requeued with --permission-mode auto (run 01M796PA86, works). Wanted: a run the daemon drives takes rhiza.daemon.permissionMode when the start names none (the daemon's own picks already do), so an owner or orchestrator can queue plan runs for accepted no-leaves change sets without choosing a mode.
