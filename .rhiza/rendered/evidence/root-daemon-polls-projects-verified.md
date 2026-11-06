# root-daemon-polls-projects-verified

evidence · `01M6SCXYRJSD6F1J1WT0XCP4EM` · actor [claude-code-cf396a2d](../actor/claude-code-cf396a2d.md)

## result

PASS

## statement

At d01fe217, go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-projects|daemon-picks|describe-golden|daemon-settings|daemon-pause|daemon-repo-session|daemon-no-root)$' passes (-v shows daemon-projects and daemon-picks ran). daemon-projects.txtar builds a root registering two initialised projects (repo, second), each with an accepted ready change set, and a bare git checkout; one `rhiza daemon --once` from the root prints two 'Started a run of', writes a start-1.log in each project's own .git/rhiza/logs, runs each start in its own Rex session (rhiza-repo, rhiza-second, none for bare), records last_started in each project's own daemon state, and leaves bare/.git/rhiza absent; daemon status prints 'Skipped project: bare: not initialised' and the JSON has project bare with reason 'not initialised'. go test ./internal/daemon/... and go test ./internal/api -run 'Daemon|Init' (incl. TestDaemonHomePrefersWorkingProject) pass. Read the diff: runDaemon loops over registered projects, skips uninitialised ones with daemon.ProjectNotInitialised, and polls each other one through daemonProj.poll with a projectCtx bound to that checkout.

## Links

### Out

- verifies [root-daemon-polls-projects](../capability/rhiza/orchestration/daemon/root-daemon-polls-projects.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
