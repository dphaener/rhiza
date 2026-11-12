# daemon-auto-update-qa-2

evidence · `01M77CR2CX94432ZSX4G8GPQ46` · actor [claude-code-af0944f4](../actor/claude-code-af0944f4.md)

## result

PASS

## statement

QA cycle 2 at 42361ffc9, built rhiza with a fake make in a scratch root whose registered project has go.mod module github.com/dphaener/rhiza and a bare origin. With autoUpdate = main, after a new commit was pushed to origin/main, rhiza daemon fast-forwarded the checkout to it and ran make update there. With off, a new origin commit was neither merged nor followed by make. With main and the module renamed, the daemon logged 'auto-update: no registered project is Rhiza's source checkout; nothing to update' and ran nothing. daemon set rhiza.daemon.autoUpdate latest exits 2 ('must be off or main'). A fresh root shows autoUpdate = off (default, machine). go test ./internal/daemon, ./internal/api -run Setting|SelfUpdate|AutoUpdate and the daemon-settings testscript pass.

## Links

### Out

- verifies [daemon-auto-update](../capability/rhiza/orchestration/daemon/daemon-auto-update.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
