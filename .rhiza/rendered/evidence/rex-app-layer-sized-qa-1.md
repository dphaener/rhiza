# rex-app-layer-sized-qa-1

evidence · `01M648A1HNTJSWDCJF1GPVGDWC` · actor [claude-code-fe283bac](../actor/claude-code-fe283bac.md)

## result

PASS

## statement

Same install: each app action opens rex.session.new_layer with bounds = cell_bounds(ctx.block_id, cols, rows) from App.Size (every app today is the 100x30 default) and runs { BIN, "tui", "app", <name> } directly, with no tui launch step; cell_bounds loads clean under the real rex config check.

## Links

### Out

- verifies [rex-app-layer-sized](../capability/rhiza/rex-app-layer-sized.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
