# quick-launch-layer-verified-2

evidence · `01M62ZXXMTN0CGD3C1FZW20PZP` · actor [claude-code-e23bd97d](../actor/claude-code-e23bd97d.md)

## result

PASS

## statement

At dd3d5b0 (attempt 2, the session.attach fix in cmd/rhiza/tui.go): built the binary and ran `rhiza rex install` into a scratch REX_CONFIG, then loaded the installed init.lua/rhiza.lua in luajit against a fake rex table: it binds ctrl+a>a to rhiza_launch, and running that action calls rex.session.new_layer with a focused block running `<bin> tui launch`. go test -count=1 ./internal/tui/... ./internal/api ./cmd/rhiza pass; with the fake Rex client, TestLayerListsEachAppByName lists each app by name and TestSelectingAnAppOpensItsModal opens one focused layer running `<bin> tui app roadmap` in the cwd and disowns the block. Read apps register through launcher.Register in readapps. The fix's diff attaches the connection with rexuikit client.Attach before surface use and passes the checked session. make fmt-check and vet pass; make self-check fails only on stale .rhiza/rendered files, not on code.

## Links

### Out

- verifies [quick-launch-layer](../capability/rhiza/quick-launch-layer.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
