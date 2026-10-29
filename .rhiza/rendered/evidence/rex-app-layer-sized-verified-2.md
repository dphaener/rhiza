# rex-app-layer-sized-verified-2

evidence · `01M6496WYQN4BJ76E2E16ZYFVJ` · actor [claude-code-38d4f6b2](../actor/claude-code-38d4f6b2.md)

## result

PASS

## statement

At 76a7035 (fix for review quality f1), in a detached checkout: go test ./internal/api -run 'TestRexCellBounds|TestRexInstall' -count=1 -v passes all 10 tests. TestRexInstallSizesEachAppLayerToItsCells registers zz-sized with Cols 77 Rows 21, runs rex.install, and finds its action's bounds = cell_bounds(ctx.block_id, 77, 21) and command { BIN, "tui", "app", "zz-sized" } with no "launch"; an app with no size gets the 100x30 default. The new TestRexCellBoundsScalesTheBlockToItsWindow ran under /opt/homebrew/bin/luajit (not skipped): cell_bounds' bounds match surface.BoundsForCells for a whole-window block, a half-width split and a 70/30 split, so from a split pane the layer is now sized from the window, not the block. The Lua reads session.view as windows[].layers[].layout{rect,block_id,split{before,after}} and kind "tiled", which matches rexuikit v0.2.0 client.LayoutNode/ViewSplit/LayerTiled, and the scale-up matches its windowCellsFrom (round(cols/rect[3])). go test ./cmd/rhiza -run TestScripts/rex-install passes (default app bounds cell_bounds(ctx.block_id, 100, 30), each app runs tui app <name>). make fmt-check vet test exits 0; make self-check fails only on stale .rhiza/rendered files from the run's model writes, nothing in code. As before, the Cols/Rows sizing is shown by the Go rex.install test, not by the txtar (review compliance f2, minor).

## Links

### Out

- verifies [rex-app-layer-sized](../capability/rhiza/rex-app-layer-sized.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
