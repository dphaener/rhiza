# self-update-settings-removed-qa-1

evidence · `01M77BC2BY9B2XBRDCGQHT9ST4` · actor [claude-code-18835d7f](../actor/claude-code-18835d7f.md)

## result

PASS

## statement

QA at ce36f205e in a scratch root: with self_update_project in rhiza.toml, or rhiza.daemon.selfUpdate under [daemon], or in git config, daemon settings exits 2 naming the key and the autoUpdate hint; daemon set rhiza.daemon.selfUpdate exits 2 with the hint. The reporting is as accepted, but it is fatal: project ls, root show and the daemon itself (daemon --once, exit 1) refuse to run on such a toml (see the QA finding).

## Links

### Out

- verifies [self-update-settings-removed](../capability/rhiza/orchestration/daemon/daemon-settings/self-update-settings-removed.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
