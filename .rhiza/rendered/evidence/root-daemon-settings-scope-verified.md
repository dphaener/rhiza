# root-daemon-settings-scope-verified

evidence · `01M6SDRNE2M388Z9JGMMVYZ8JE` · actor [claude-code-35df0f99](../actor/claude-code-35df0f99.md)

## result

PASS

## statement

Verifier at a67915229a: go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/(daemon-settings-scope|daemon-settings|usage-settings|describe-golden|run-ship-change|ship-deploy-settings|stage-autonomy) all PASS. daemon-settings-scope.txtar sets rhiza.daemon.maxRuns (found in $WORK/rhiza.toml, absent from git config) and rhiza.autonomy.ship (found in project a's git config), sees both in daemon settings with (rhiza.toml, machine)/(git config, project) and in --json with scope and source, sees rhiza.daemon.trimGoCache with machine scope, and fails a project key set at the root without --project; with two projects each with its own rhiza.eventCommand, one daemon --once poll writes A's two moves only to out-a and B's only to out-b. By hand with a scratch build and root: daemon set rhiza.autonomy.ship owner at the root exits 2 naming --project; daemon settings --project a --json lists trimGoCache scope machine. internal/rootcfg tests pass; internal/api TestDaemon*, TestShipDeploy*, TestStatusHook*, daemon_projects_test and usage_hold_test (27 tests) pass. The full internal/api package timed out on this loaded host (load 8 on 4 cores) in the read-path budget tests (TestReadBudget, TestInboxGateQuestionsBudget), which this diff does not touch; CI runs those. Diff read: settingDefs gain a project flag, machine keys read/write rhiza.toml [daemon] via rootcfg.SetDaemon, project keys git config via scopeCtx/projectCtx, hooks resolved per project in daemon_projects.go.

## Links

### Out

- verifies [root-daemon-settings-scope](../capability/rhiza/orchestration/daemon/root-daemon-settings-scope.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
