# self-update-settings-removed-qa-2

evidence · `01M77CR26TSWAGJ0PQ8RKTTBZQ` · actor [claude-code-af0944f4](../actor/claude-code-af0944f4.md)

## result

PASS

## statement

QA cycle 2 at 42361ffc9, built rhiza in a scratch root (own HOME, fake make on PATH) whose rhiza.toml has the VM's shape: self_update_project = "rhiza" and [daemon] rhiza.daemon.selfUpdate = "true". project ls, daemon settings and health exit 0; daemon settings shows rhiza.daemon.autoUpdate = main (rhiza.toml, machine). rhiza daemon starts, logs both old keys as removed (selfUpdate read as autoUpdate = main), and after origin/main moved it fast-forwarded the checkout and ran make update. rhiza daemon set rhiza.daemon.permissionMode acceptEdits rewrote rhiza.toml without either old key and with autoUpdate = main, and wrote no git config. A git config rhiza.daemon.selfUpdate is logged as removed and ignored at daemon start. go test ./internal/rootcfg passes.

## Links

### Out

- verifies [self-update-settings-removed](../capability/rhiza/orchestration/daemon/daemon-settings/self-update-settings-removed.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
