# rex-key-mode-verified-2

evidence · `01M649C1EC0ZHPKMVCK67HDCEW` · actor [claude-code-f0becb64](../actor/claude-code-f0becb64.md)

## result

PASS

## statement

Verified at ff7cb93 (attempt 2). go test ./cmd/rhiza -run TestScripts/rex-install passes: the rex-install script checks the launch key enters the rhiza mode (client.mode.enter), default keys for roadmap and model, the keys = {} settings map, session-ls mapped to rhiza_app_session_ls in APP_ACTIONS, and rex.bind(MODE .. "/" .. k, APP_ACTIONS[app]). All TestRexInstall* Go tests pass, including BindsAKeyPerAppInTheKeyMode, which now requires every app's bound action to be one a rex.action defines, and KeepsKeyOverridesOnReinstall. I built the binary, ran rex install against a scratch REX_CONFIG with a fake rex, and ran the written rhiza.lua under luajit with a stub rex. It made 16 rhiza/<key> binds, one per registered app, with 16 distinct keys, each to a rhiza_app_* action the same file defines (this includes the 7 hyphenated apps, the blocker from qa f2). ctrl+a>a enters the mode. Setting keys = { roadmap = "x" } rebinds roadmap to rhiza/x. This matches the owner's decision rex-app-keys-default-overridable. make check: fmt, vet and tests pass; self-check fails only because .rhiza/rendered actor files and index.md are stale from this run's model writes, not because of this leaf's code. Still open: the unconfirmed minor quality f3, where a remap onto another app's default key collides.

## Links

### Out

- verifies [rex-key-mode](../capability/rhiza/rex-key-mode.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
