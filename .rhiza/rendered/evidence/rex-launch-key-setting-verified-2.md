# rex-launch-key-setting-verified-2

evidence · `01M649K95RQ8ZJYQPC6R0K00EG` · actor [claude-code-34dde8a0](../actor/claude-code-34dde8a0.md)

## result

PASS

## statement

At 9281614 (attempt 2, after the fix that revised quick-launch-layer to the palette-only rhiza_launch action), go test ./cmd/rhiza -run TestScripts/rex-install -count=1 passes. The script asserts a first install writes launch_key = "ctrl+a>a" (internal/api/rex_install.go:41, the default in the setting's > spacing), sed-edits it to ctrl+b>x, reinstalls with --json reporting written [] and cmp showing rhiza.lua unchanged, the edited line still present exactly once, and the -- settings begin / -- settings end block kept. go test ./internal/api -run Rex and go test ./... pass. Read the diff: the fix commit only updates quick-launch-layer's statement and acceptance in .rhiza and records its actor; the launch_key code is unchanged from the verified attempt 1. make check's self-check fails only on rendered actor pages for sessions recorded after this commit (claude-code-34dde8a0, claude-code-4dee40ee) and index.md, the same at the parent commit, not on anything this leaf wrote.

## Links

### Out

- verifies [rex-launch-key-setting](../capability/rhiza/rex-launch-key-setting.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
