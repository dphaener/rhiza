# self-update-settings-removed-verified-2

evidence · `01M77CA1717GQ8F504WR432MPM` · actor [claude-code-be281453](../actor/claude-code-be281453.md)

## result

PASS

## statement

At 14ad701b8 I read the diff (rootcfg.Load carries self_update_project and rhiza.daemon.selfUpdate over into Config.Migrated notes naming rhiza.daemon.autoUpdate, selfUpdate = true read as autoUpdate = main unless autoUpdate is set; runDaemon prints the notes; SetDaemon rewrites without the old keys; daemon set rhiza.daemon.selfUpdate refused). go test -trimpath ./internal/rootcfg/... passed (TestLoadMigratesRemovedSelfUpdate: four tomls, autoUpdate main/empty/empty/off, notes name rhiza.daemon.autoUpdate, a write drops the old keys and reloads with no notes). go test -trimpath ./internal/api -run 'Daemon|Setting' -count=1 from a git archive copy outside ~/code passed, incl. TestDaemonSettingsMigrateRemovedSelfUpdate (settings read autoUpdate = main, two notes naming rhiza.daemon.autoUpdate, daemon set maxRuns drops selfUpdate and self_update_project and keeps autoUpdate, daemon set selfUpdate refused naming autoUpdate) and TestDaemonSetBrokenRootWritesNoGitConfig. Run in place under ~/code, the same tests passed but the package's root-state guard flagged ~/code/.rhiza-root/usage.live.json changing, which is the live root daemon's file, not this change.

## Links

### Out

- verifies [self-update-settings-removed](../capability/rhiza/orchestration/daemon/daemon-settings/self-update-settings-removed.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
