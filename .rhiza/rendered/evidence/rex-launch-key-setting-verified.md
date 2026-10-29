# rex-launch-key-setting-verified

evidence · `01M6486XAJ5RTVW7YPNWF7VC92` · actor [claude-code-a11b0513](../actor/claude-code-a11b0513.md)

## result

PASS

## statement

At 2b1b299 ran go test ./cmd/rhiza -run TestScripts/rex-install -count=1 and go test ./internal/api -run TestRexInstall -count=1: both PASS. The rex-install script checks a first install writes launch_key = "ctrl+a>a" in conf/rhiza.lua (rexDefaultSettings in internal/api/rex_install.go now writes ctrl+a>a, changed with the test), then seds launch_key to "ctrl+b>x", reinstalls with --json (written: []), cmps rhiza.lua byte for byte with the edited copy, greps the edited line once, and greps the -- settings begin / -- settings end markers. The generated Lua binds settings.launch_key (spaces turned to >) to client.mode.enter. make fmt-check vet test exit 0; make self-check fails only on stale .rhiza/rendered actor files and index from later sessions, unrelated to this leaf.

## Links

### Out

- verifies [rex-launch-key-setting](../capability/rhiza/rex-launch-key-setting.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
