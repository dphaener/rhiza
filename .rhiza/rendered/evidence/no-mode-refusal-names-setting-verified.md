# no-mode-refusal-names-setting-verified

evidence · `01M7995DHF02QZC0R9640EQ5X8` · actor [claude-code-8035eb5d](../actor/claude-code-8035eb5d.md)

## result

PASS

## statement

Verifier at 240df0cec: read the diff (internal/api/run.go checkAgentStep message now 'give --permission-mode, set rhiza.daemon.permissionMode (for a run the daemon starts), or set permission_mode on the step'; run_test.go NoMode case checks that text; run-start-queue.txtar adds the unset and --foreground cases). Ran go test -trimpath ./internal/api -run 'TestRunStartAgentRefusalsBeforeWrite|TestRunStart|TestQueuedCheckMode' -count=1: ok. Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-start-queue|daemon-picks)$' -count=1 -v: both PASS. In run-start-queue, with rhiza.daemon.permissionMode unset, 'rhiza run start --workflow plan' fails with stderr naming rhiza.daemon.permissionMode and starts.json holds 0 change entries before and after; with the setting = manual, 'rhiza run start --workflow plan --foreground' exits 2 with stderr: step "interview" is an agent step with no permission mode: give --permission-mode, set rhiza.daemon.permissionMode (for a run the daemon starts), or set permission_mode on the step.

## Links

### Out

- verifies [no-mode-refusal-names-setting](../capability/rhiza/orchestration/daemon/run-start-through-daemon/no-mode-refusal-names-setting.md)

## History

- `01M797TRSFCEHV3FXTYVRV7D74` Fix A run queued by hand (rhiza run start <cs> --workflow plan --autonomous) does not get the daemon's own permission mode: the CLI refuses --permission-mode bypassPermissions, and with none the builtin plan's interview step is refused ('agent step with no permission mode'), so a queued plan run must pick a weaker mode. On 2026-11-13 ~00:40 plan run 01M796JBG1 (change set 01M76XR5, queued with --permission-mode acceptEdits) stopped at its first shell command on a permission prompt (inbox 84109a8). rhiza run abandon then refused it ('running in another process'), so the orchestrator killed its driver and agent by PID, abandoned it, and requeued with --permission-mode auto (run 01M796PA86, works). Wanted: a run the daemon drives takes rhiza.daemon.permissionMode when the start names none (the daemon's own picks already do), so an owner or orchestrator can queue plan runs for accepted no-leaves change sets without choosing a mode.
