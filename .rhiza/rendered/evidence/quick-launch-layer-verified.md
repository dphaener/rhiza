# quick-launch-layer-verified

evidence · `01M62X1H7QFTYJHE0AB1VNKV3D` · actor [claude-code-dd70ce51](../actor/claude-code-dd70ce51.md)

## result

PASS

## statement

Verified at 2dc1ec4. go test -count=1 ./internal/tui/launcher ./internal/api: TestLayerListsEachAppByName (each app name in the ui.List view), TestSelectingAnAppOpensItsModal (with a fake surface.Rex, down+enter opens one focused layer whose layout runs /opt/rhiza/bin/rhiza tui app roadmap in /work, plus disown_block), TestEscapeClosesWithoutOpening, TestAppsAreSortedAndReplaced (registry) and TestRexInstallScriptOpensTheLauncher (launch_key ctrl+a a, BIN, rhiza_launch action and bind; key kept, BIN refreshed on reinstall) all PASS. Also built rhiza, ran rex install with REX_CONFIG in a scratch dir, and loaded the written rhiza.lua in luajit against a fake rex: it binds ctrl+a>a to rhiza_launch, whose run opens a focused new_layer with block flavor com.superlogical.terminal.shell running <BIN> tui launch. The Lua calls match rexuikit v0.2.0's Rex API docs. make check (RHIZA_ROOT unset) exits 0. Note: rhiza tui app <name> and the registered apps come from the read-apps leaf; this leaf accepts only tui launch.

## Links

### Out

- verifies [quick-launch-layer](../capability/rhiza/quick-launch-layer.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
