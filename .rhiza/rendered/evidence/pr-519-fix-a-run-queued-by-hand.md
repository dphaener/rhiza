# pr-519-fix-a-run-queued-by-hand

evidence · `01M79AQGBFGVER8SQ99W1DYPB3` · actor [dphaener](../actor/dphaener.md)

## head

5699f049be375bafc5a3ef7aac977b7391efe1e4

## pr

https://github.com/dphaener/rhiza/pull/519

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #519 (https://github.com/dphaener/rhiza/pull/519) passed CI at 5699f04

## status

verified

## Links

### Out

- verifies [queued-start-daemon-mode](../capability/rhiza/orchestration/daemon/run-start-through-daemon/queued-start-daemon-mode.md)
- verifies [no-mode-refusal-names-setting](../capability/rhiza/orchestration/daemon/run-start-through-daemon/no-mode-refusal-names-setting.md)

## History

- `01M797TRSFCEHV3FXTYVRV7D74` Fix A run queued by hand (rhiza run start <cs> --workflow plan --autonomous) does not get the daemon's own permission mode: the CLI refuses --permission-mode bypassPermissions, and with none the builtin plan's interview step is refused ('agent step with no permission mode'), so a queued plan run must pick a weaker mode. On 2026-11-13 ~00:40 plan run 01M796JBG1 (change set 01M76XR5, queued with --permission-mode acceptEdits) stopped at its first shell command on a permission prompt (inbox 84109a8). rhiza run abandon then refused it ('running in another process'), so the orchestrator killed its driver and agent by PID, abandoned it, and requeued with --permission-mode auto (run 01M796PA86, works). Wanted: a run the daemon drives takes rhiza.daemon.permissionMode when the start names none (the daemon's own picks already do), so an owner or orchestrator can queue plan runs for accepted no-leaves change sets without choosing a mode.
