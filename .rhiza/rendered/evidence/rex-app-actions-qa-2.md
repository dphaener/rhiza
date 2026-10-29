# rex-app-actions-qa-2

evidence · `01M64A7SN9GB63YFN68SCHNXJS` · actor [claude-code-afc79974](../actor/claude-code-afc79974.md)

## result

PASS

## statement

At 7b2fd58, built go build -o $TMPDIR/rhiza-qa2 ./cmd/rhiza and ran rhiza rex install --json with REX_CONFIG pointing at a scratch init.lua: rhiza.lua holds one rex.action per registered app (16, rhiza_app_<name>, titled with the app's name, command { BIN, "tui", "app", <name> }) and no rhiza_launch. Real Rex's rex config check on it reports actions: 16, modes: 1, bindings: 18, no diagnostics. Against the live Rex server, rex do --action rhiza_app_ready / rhiza_app_model / rhiza_app_roadmap <rhiza.lua> each opened a layer whose block ran rhiza-qa2 tui app <name> in the invoking block's cwd; rex capture showed the ready, model (tree) and roadmap screens; q quit and closed the layer. rhiza tui launch now exits 2 with the tui app usage; go test ./internal/api -run Rex (incl. TestRexInstallRegistersAnActionPerApp) and the rex-install script pass.

## Links

### Out

- verifies [rex-app-actions](../capability/rhiza/rex-app-actions.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
