# root-daemon-service-verified-2

evidence · `01M6SJDC20HD578TBM5TGAD5EX` · actor [claude-code-ab28259b](../actor/claude-code-ab28259b.md)

## result

PASS

## statement

Verified at 4eb72fe5 (attempt 2). go test -trimpath ./cmd/rhiza -run TestScripts/(daemon-install|daemon-start-stop|daemon-projects|daemon-status|root-init) all PASS: daemon-install, with fake systemctl/loginctl on PATH, runs daemon install from a subfolder of a project and finds one unit rhiza-root.service with WorkingDirectory=$WORK (the root, not the repo), then daemon start/stop from the root call systemctl --user start/stop rhiza-root.service and uninstall calls disable --now plus daemon-reload (cmp against want logs). The fix commit makes daemon status resolve the root (rootcfg.Find) and read the root lock plus the home project's state recorded in .rhiza-root/daemon.home; daemon-projects now checks status from the root, outside any checkout, and from a non-home project gives the same JSON, and a live root daemon reports running from both. By hand with a scratch build: rhiza root init . then rhiza daemon status at the fresh root printed 'Daemon not running' and exited 0 (earlier QA failure was exit 1 'sessions need a git repository'). go test ./internal/daemon/... ok; go vet and gofmt clean on internal/api and internal/daemon.

## Links

### Out

- verifies [root-daemon-service](../capability/rhiza/orchestration/daemon/root-daemon-service.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
