# daemon-auto-update-verified

evidence · `01M77APJY7ZGF9SRQYGY8QXH1D` · actor [claude-code-04ee5483](../actor/claude-code-04ee5483.md)

## result

PASS

## statement

At 26fb806ba I read the diff and ran the scoped tests in a detached worktree outside ~/code with fakes only. TestDaemonAutoUpdateSetting (internal/api; fake git and make) passes: autoUpdate unset or off runs no git or make; main with a registered project whose go.mod is not Rhiza's logs 'no registered project is Rhiza's source checkout' and runs nothing; main with a Rhiza checkout runs git fetch origin main, git merge --ff-only and make update exactly once over two polls. internal/daemon passes, including TestRhizaCheckoutFindsModule (finds the project by module path, false when there is none) and the existing SelfUpdate tests. TestScripts/daemon-settings passes and shows rhiza.daemon.autoUpdate = off (default, machine). gofmt and go vet are clean. Inside ~/code the api package fails only on the enclosing-root guard, because the live daemon rewrites .rhiza-root/usage.live.json; the test does not cause it.

## Links

### Out

- verifies [daemon-auto-update](../capability/rhiza/orchestration/daemon/daemon-auto-update.md)

## History

- `01M763Z16T0ZGQMY1EC2Y81V0K` Auto-update setting: rhiza.daemon.selfUpdate becomes rhiza.daemon.autoUpdate with values off (default) and main (follow a registered Rhiza source checkout and run make update, as self-update does today); rhiza.toml's self_update_project is removed; the setting is in daemon settings, the generated reference and the user docs.
