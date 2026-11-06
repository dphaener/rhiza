# root-daemon-service-verified

evidence · `01M6SD8ZTCNX3VGC5PRSS02905` · actor [claude-code-b2f10654](../actor/claude-code-b2f10654.md)

## result

PASS

## statement

Verifier at 29051fdd: ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-install|daemon-start-stop|describe-golden)$' -count=1 -v; all three PASS. daemon-install.txtar puts fake systemctl/loginctl on PATH, runs rhiza daemon install from repo/sub (inside a project under the root $WORK, which holds rhiza.toml), and checks the one unit rhiza-root.service has WorkingDirectory=$WORK (not .../repo), ExecStart ... daemon --autonomous, Restart=always, After=rex.service, and that systemctl --user enable --now rhiza-root.service ran; then from $WORK itself (no git checkout) rhiza daemon start and stop log systemctl --user start/stop rhiza-root.service, and uninstall logs systemctl --user disable --now rhiza-root.service plus daemon-reload and removes the unit. Diff read: serviceSetup resolves the root with rootcfg.Find (usage error with next rhiza root init outside a root) instead of git worktree list; daemonMode keys on the fixed unit name; start/stop in service mode no longer require git.

## Links

### Out

- verifies [root-daemon-service](../capability/rhiza/orchestration/daemon/root-daemon-service.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
