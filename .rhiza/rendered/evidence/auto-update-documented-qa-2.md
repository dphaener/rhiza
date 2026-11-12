# auto-update-documented-qa-2

evidence · `01M77CR2JK7YBMYTSEMNGNM9JK` · actor [claude-code-af0944f4](../actor/claude-code-af0944f4.md)

## result

PASS

## statement

QA cycle 2 at 42361ffc9: the built rhiza docs --out into a copy of docs/reference gives no diff; config.md has rhiza.daemon.autoUpdate; docs/runs-and-daemon.md names autoUpdate with off and main; no selfUpdate or self_update_project left in docs/runs-and-daemon.md, docs/reference/config.md or README.md. go test -trimpath ./internal/doclinks passes.

## Links

### Out

- verifies [auto-update-documented](../capability/rhiza/orchestration/daemon/daemon-settings/auto-update-documented.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
