# queued-run-mode-interview

source · `01M7982437J52G5PWT1THXVQCS`

## quotes

```json
[]
```

## raw_sha256

e8d794ab6c0ce726b122107598b44b7f2502892b648ad5973016885323fe349b

## session

16742c4a-fa0f-4ac5-b6a2-01a1d27db9cb

## summary

Autonomous plan run of bug change set 01M797T (no owner). The planner read run.go (checkAgentSteps refuses a queued start with no --permission-mode at queue time), queueRunStart (records the start's mode, possibly empty) and daemon.go's queued-start loop (fills an empty mode with rhiza.daemon.permissionMode). On this dedicated VM a dry-run with bypassPermissions passes the trust check; with no mode it is refused ('agent step with no permission mode'). The planner asked two questions (queue-time check, and the no-mode case); the arbiter answered both with Decision queued-run-daemon-mode: check at queue time with the daemon's setting standing in, queue with no mode so the daemon applies its setting at start; keep the refusal when the setting is unset or the start is --foreground, with the message naming rhiza.daemon.permissionMode. The change set serves m4-orchestrated-by-rhiza: the incident needed an orchestrator to kill processes by PID. Leaves are under run-start-through-daemon. Not breaking (only a refusal becomes a queue), not severe.

## title

Plan interview: a queued run takes the daemon's permission mode

## Links

### In

- derived-from [queued-run-daemon-mode](../decision/queued-run-daemon-mode.md)

## History

- `01M797TRSFCEHV3FXTYVRV7D74` Fix A run queued by hand (rhiza run start <cs> --workflow plan --autonomous) does not get the daemon's own permission mode: the CLI refuses --permission-mode bypassPermissions, and with none the builtin plan's interview step is refused ('agent step with no permission mode'), so a queued plan run must pick a weaker mode. On 2026-11-13 ~00:40 plan run 01M796JBG1 (change set 01M76XR5, queued with --permission-mode acceptEdits) stopped at its first shell command on a permission prompt (inbox 84109a8). rhiza run abandon then refused it ('running in another process'), so the orchestrator killed its driver and agent by PID, abandoned it, and requeued with --permission-mode auto (run 01M796PA86, works). Wanted: a run the daemon drives takes rhiza.daemon.permissionMode when the start names none (the daemon's own picks already do), so an owner or orchestrator can queue plan runs for accepted no-leaves change sets without choosing a mode.
