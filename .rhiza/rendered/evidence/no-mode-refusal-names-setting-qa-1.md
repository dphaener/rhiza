# no-mode-refusal-names-setting-qa-1

evidence · `01M799H7AKATM172NN1D4HB2NJ` · actor [claude-code-82543a16](../actor/claude-code-82543a16.md)

## result

PASS

## statement

QA at 30e69b6fe with a binary built from the branch, in a scratch Rhiza root and repo. With rhiza.daemon.permissionMode unset, rhiza run start --workflow plan (no --permission-mode) exits 2 with: step "interview" is an agent step with no permission mode: give --permission-mode, set rhiza.daemon.permissionMode (for a run the daemon starts), or set permission_mode on the step; .git/rhiza does not exist afterwards, so nothing was queued. With the setting = manual, rhiza run start --workflow plan --foreground exits 2 with the same message.

## Links

### Out

- verifies [no-mode-refusal-names-setting](../capability/rhiza/orchestration/daemon/run-start-through-daemon/no-mode-refusal-names-setting.md)

## History

- `01M797TRSFCEHV3FXTYVRV7D74` Fix A run queued by hand (rhiza run start <cs> --workflow plan --autonomous) does not get the daemon's own permission mode: the CLI refuses --permission-mode bypassPermissions, and with none the builtin plan's interview step is refused ('agent step with no permission mode'), so a queued plan run must pick a weaker mode. On 2026-11-13 ~00:40 plan run 01M796JBG1 (change set 01M76XR5, queued with --permission-mode acceptEdits) stopped at its first shell command on a permission prompt (inbox 84109a8). rhiza run abandon then refused it ('running in another process'), so the orchestrator killed its driver and agent by PID, abandoned it, and requeued with --permission-mode auto (run 01M796PA86, works). Wanted: a run the daemon drives takes rhiza.daemon.permissionMode when the start names none (the daemon's own picks already do), so an owner or orchestrator can queue plan runs for accepted no-leaves change sets without choosing a mode.
