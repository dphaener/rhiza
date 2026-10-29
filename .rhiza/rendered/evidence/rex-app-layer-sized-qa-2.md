# rex-app-layer-sized-qa-2

evidence · `01M64A7T9R4VY5MQW2NDMG67XR` · actor [claude-code-afc79974](../actor/claude-code-afc79974.md)

## result

PASS

## statement

At 7b2fd58, each app action in the installed rhiza.lua calls cell_bounds(ctx.block_id, 100, 30) (the App.Size defaults) and runs tui app <name> directly, with no launch step. Live via rex do --action: from a single tiled block (80x24, rect 1x1) the model/roadmap layers got bounds {0,0,1,1}, the app clamped to the whole window as the model says when the window is smaller than the app; from my own block, nested 18 splits deep (rect width 1.5e-05, size 80x24 as no client renders the session), the ready layer got bounds {0.49999, 0.171, 2.0e-05, 0.657}: height matches the formula for 48 window rows, width is a sliver because window_cells scales 80 cols by 1/1.5e-05. TestRexCellBoundsScalesTheBlockToItsWindow and TestRexInstallSizesEachAppLayerToItsCells pass. Sizing works for ordinary layouts; the degenerate case is reported as a minor finding.

## Links

### Out

- verifies [rex-app-layer-sized](../capability/rhiza/rex-app-layer-sized.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
