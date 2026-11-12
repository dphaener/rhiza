# auto-update-documented-qa-1

evidence · `01M77BC29T1M645H5HECSBGF7V` · actor [claude-code-18835d7f](../actor/claude-code-18835d7f.md)

## result

PASS

## statement

QA at ce36f205e: the built rhiza docs --out into a copy of docs/reference gives no diff; config.md has rhiza.daemon.autoUpdate; docs/runs-and-daemon.md names autoUpdate with off and main and latest/release as coming with release channels; no selfUpdate or self_update_project left in docs/ or README.md. go test -trimpath ./internal/doclinks passes.

## Links

### Out

- verifies [auto-update-documented](../capability/rhiza/orchestration/daemon/daemon-settings/auto-update-documented.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
