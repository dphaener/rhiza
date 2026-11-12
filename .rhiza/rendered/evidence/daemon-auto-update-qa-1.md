# daemon-auto-update-qa-1

evidence · `01M77BC26QN0WH9YNF2GFV3HEE` · actor [claude-code-18835d7f](../actor/claude-code-18835d7f.md)

## result

PASS

## statement

QA at ce36f205e with a -trimpath build in a scratch root (HOME and root outside ~/code, port off): rhiza daemon settings shows rhiza.daemon.autoUpdate = off (default, machine); daemon set takes off and main and refuses latest and true (exit 2). Ran the real daemon loop (timeout -s INT 25) with a registered project whose go.mod declares github.com/dphaener/rhiza, origin/main one commit ahead and a Makefile whose update target touches a marker: with main it fast-forwarded the checkout and logged 'self-update: updated to 52b44f898' and the marker appeared; with off and origin ahead again, nothing ran; with main and that project unregistered it logged 'auto-update: no registered project is Rhiza's source checkout; nothing to update' once and ran no make.

## Links

### Out

- verifies [daemon-auto-update](../capability/rhiza/orchestration/daemon/daemon-auto-update.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
