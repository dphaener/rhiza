# rex-app-actions-verified

evidence · `01M6478KZ0VCK4YT5C9XW98BTK` · actor [claude-code-3b240678](../actor/claude-code-3b240678.md)

## result

PASS

## statement

At 6acb2fe: go test ./cmd/rhiza -run TestScripts/rex-install passes; the script test greps rhiza.lua for rex.action entries titled with app names from readapps (roadmap, session-ls), writeapps (accept) and modelapp (model), each running "tui", "app", <name>. go test ./internal/api -run TestRexInstall -count=2 passes; TestRexInstallRegistersAnActionPerApp registers zz-extra, installs and finds its action, checks for 1+len(launcher.Apps()) rex.action calls, then registers zz-later, reinstalls and finds its action. A scratch build with a fake rex and REX_CONFIG wrote rhiza.lua with 17 rex.action blocks: rhiza_launch plus one for each of the 16 registered apps (7 readapps, 8 writeapps, model). Each block is titled with the app's name and runs {BIN, "tui", "app", <name>} in a layer. A reinstall was byte-identical. In make check, fmt, vet and every go test pass; self-check reports only stale .rhiza/rendered actor and index files from run state, which this leaf's diff does not touch.

## Links

### Out

- verifies [rex-app-actions](../capability/rhiza/rex-app-actions.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
