# self-update-settings-removed-verified

evidence · `01M779WKB46T58NZDW8F4NNQFG` · actor [claude-code-445deb58](../actor/claude-code-445deb58.md)

## result

PASS

## statement

Verifier claude-code-445deb58 at a05cacac6: go test -trimpath ./internal/rootcfg/... ok (TestLoadRejectsRemovedSelfUpdate: rhiza.toml with self_update_project, and with [daemon] rhiza.daemon.selfUpdate, each makes Load return ErrRemoved with a message naming rhiza.daemon.autoUpdate); go test -trimpath ./internal/daemon/... ok; go test -trimpath ./internal/api -run TestDaemonSettingsReportRemovedSelfUpdate -v PASS (daemon settings reports both removed rhiza.toml keys with the autoUpdate hint; rhiza daemon set rhiza.daemon.selfUpdate true fails naming rhiza.daemon.autoUpdate); go test ./internal/api -run 'Daemon|Setting|SelfUpdate': every test passed, the package then failing only on the root-state guard (.rhiza-root/usage.live.json changed by the live daemon under ~/code, not by this diff); go vet clean. Diff read: selfUpdate setting def, SelfUpdateProject field, SelfUpdateCheckout/ForCheckout and the daemon's self-update wiring are removed; describe-golden, cli.md and config.md updated; git config rhiza.daemon.selfUpdate is also reported. Note: Load now errors on these keys, and the live root ~/code/rhiza.toml still sets both self_update_project and rhiza.daemon.selfUpdate, so those lines must be removed when this lands.

## Links

### Out

- verifies [self-update-settings-removed](../capability/rhiza/orchestration/daemon/daemon-settings/self-update-settings-removed.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
