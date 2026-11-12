# daemon-auto-update-verified-2

evidence · `01M77CB2H6E56MB9WXG7F53WVW` · actor [claude-code-166839ca](../actor/claude-code-166839ca.md)

## result

PASS

## statement

Verifier, attempt 2, at cd74e6466: read the diff (internal/daemon/selfupdate.go RhizaCheckout by go.mod module path, internal/api/daemon.go autoUpdate wiring with the once-only none-registered log, daemon_settings.go autoUpdate off|main default off). Ran go test -trimpath -count=1 ./internal/daemon/... ./internal/rootcfg/... (ok); go test ./internal/api -run 'AutoUpdate|SelfUpdate|RhizaCheckout|DaemonSet|Settings' -v: TestDaemonAutoUpdateSetting PASS, with fake git and make covering unset and off (nothing run), main with a Rhiza checkout (git fetch origin main, git merge --ff-only, make update once over two ticks) and main with none (logged, nothing run); TestRhizaCheckoutFindsModule passes in internal/daemon. go test ./cmd/rhiza -run 'TestScripts/(daemon-settings|describe-golden)$' ok; daemon-settings asserts rhiza.daemon.autoUpdate = off (default, machine). The off/main/none cases live in internal/api rather than internal/daemon, as the compliance check noted (minor); the behaviour the acceptance names is covered.

## Links

### Out

- verifies [daemon-auto-update](../capability/rhiza/orchestration/daemon/daemon-auto-update.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
