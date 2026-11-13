# no-mode-refusal-names-setting

capability · `01M7982F49NA91KNX67QHP9914` · parent [run-start-through-daemon](../run-start-through-daemon.md)

## acceptance

A testscript: with rhiza.daemon.permissionMode unset, rhiza run start <cs> --workflow plan with no --permission-mode fails, its stderr names rhiza.daemon.permissionMode, and nothing is queued; with the setting set, rhiza run start --foreground with no --permission-mode is still refused; the NoMode case in internal/api/run_test.go checks the new message

## statement

A run start that names no --permission-mode is still refused when rhiza.daemon.permissionMode is unset or the start is --foreground, and the refusal names rhiza.daemon.permissionMode as a way out next to --permission-mode and the step's permission_mode

## Links

### Out

- decided-by [queued-run-daemon-mode](../../../../../decision/queued-run-daemon-mode.md)

### In

- verifies [no-mode-refusal-names-setting-verified](../../../../../evidence/no-mode-refusal-names-setting-verified.md)
- verifies [no-mode-refusal-names-setting-qa-1](../../../../../evidence/no-mode-refusal-names-setting-qa-1.md)
- verifies [pr-519-fix-a-run-queued-by-hand](../../../../../evidence/pr-519-fix-a-run-queued-by-hand.md)

## History

- `01M797TRSFCEHV3FXTYVRV7D74` Fix A run queued by hand (rhiza run start <cs> --workflow plan --autonomous) does not get the daemon's own permission mode: the CLI refuses --permission-mode bypassPermissions, and with none the builtin plan's interview step is refused ('agent step with no permission mode'), so a queued plan run must pick a weaker mode. On 2026-11-13 ~00:40 plan run 01M796JBG1 (change set 01M76XR5, queued with --permission-mode acceptEdits) stopped at its first shell command on a permission prompt (inbox 84109a8). rhiza run abandon then refused it ('running in another process'), so the orchestrator killed its driver and agent by PID, abandoned it, and requeued with --permission-mode auto (run 01M796PA86, works). Wanted: a run the daemon drives takes rhiza.daemon.permissionMode when the start names none (the daemon's own picks already do), so an owner or orchestrator can queue plan runs for accepted no-leaves change sets without choosing a mode.
