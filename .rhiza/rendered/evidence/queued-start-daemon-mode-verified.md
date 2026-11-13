# queued-start-daemon-mode-verified

evidence · `01M7998QSGV3EA5G65FJPYSN3S` · actor [claude-code-0b7369d7](../actor/claude-code-0b7369d7.md)

## result

PASS

## statement

Verified at 401ac36c0 as claude-code-0b7369d7. Read the diff: run.go adds queuedCheckMode (the start's own mode, else rhiza.daemon.permissionMode for a queued start, none for --foreground) and passes it to both checkAgentSteps calls in runRunStart; queueRunStart still records r.PermissionMode, so the queued start holds no mode, and daemonQueuedStarts fills an empty mode with the daemon's own. go test -trimpath ./internal/api/ -run 'TestQueuedCheckMode|TestRunStart|TestDaemon': ok (TestQueuedCheckModeStandsInDaemonMode covers Queued->manual, Named->plan, Foreground->empty; TestRunStartAgentRefusalsBeforeWrite still passes). go test -trimpath ./cmd/rhiza -run 'TestScripts/run-start-queue$' -testwork: PASS. In its log, with rhiza.daemon.permissionMode=nonsense, run start --workflow plan is refused with 'claude-code: unknown permission mode "nonsense"' (exit 2) and nothing is queued (the next start on the same change set queues without 'already has a run queued'); with it set to manual, run start --workflow plan --autonomous queues, starts.json has 0 permission_mode keys, and rhiza daemon --max-runs 3 --once starts it. I opened the kept work dir: the daemon-started run's journal for that change set records permission_mode manual and autonomous true, so the daemon passed its own setting. Minor: the script does not assert that journal mode itself; I checked it by hand.

## Links

### Out

- verifies [queued-start-daemon-mode](../capability/rhiza/orchestration/daemon/run-start-through-daemon/queued-start-daemon-mode.md)

## History

- `01M797TRSFCEHV3FXTYVRV7D74` Fix A run queued by hand (rhiza run start <cs> --workflow plan --autonomous) does not get the daemon's own permission mode: the CLI refuses --permission-mode bypassPermissions, and with none the builtin plan's interview step is refused ('agent step with no permission mode'), so a queued plan run must pick a weaker mode. On 2026-11-13 ~00:40 plan run 01M796JBG1 (change set 01M76XR5, queued with --permission-mode acceptEdits) stopped at its first shell command on a permission prompt (inbox 84109a8). rhiza run abandon then refused it ('running in another process'), so the orchestrator killed its driver and agent by PID, abandoned it, and requeued with --permission-mode auto (run 01M796PA86, works). Wanted: a run the daemon drives takes rhiza.daemon.permissionMode when the start names none (the daemon's own picks already do), so an owner or orchestrator can queue plan runs for accepted no-leaves change sets without choosing a mode.
