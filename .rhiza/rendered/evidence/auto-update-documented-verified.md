# auto-update-documented-verified

evidence · `01M77A2E93DNPXGHPZE297ADB0` · actor [claude-code-a3026434](../actor/claude-code-a3026434.md)

## result

PASS

## statement

At 328e4c60f I built rhiza with -trimpath and ran rhiza docs --out into a copy of docs/reference: no diff from the committed docs/reference, and config.md shows rhiza.daemon.autoUpdate with type off or main, default off, and latest and release named as coming with the release channels milestone. go test -trimpath ./internal/doclinks -run TestDocsDescribeAutoUpdate passes (TestRepoDocsLinksResolve too). In a git archive copy, renaming rhiza.daemon.autoUpdate in docs/runs-and-daemon.md made it fail, and adding selfUpdate to it made it fail. docs/runs-and-daemon.md names autoUpdate with off and main and no longer names selfUpdate or self_update_project. go test ./internal/docs and ./internal/api -run 'Setting|Docs' pass.

## Links

### Out

- verifies [auto-update-documented](../capability/rhiza/orchestration/daemon/daemon-settings/auto-update-documented.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
