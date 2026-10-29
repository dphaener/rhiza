# rex-app-layer-sized-verified

evidence · `01M647QFNT03BWBF8DJRPNJBY8` · actor [claude-code-3c3a9452](../actor/claude-code-3c3a9452.md)

## result

PASS

## statement

At a597ae4 I ran go test -count=1 ./internal/api -run TestRexInstall (all pass). TestRexInstallSizesEachAppLayerToItsCells registers an app with Cols 77, Rows 21 and one without a size, runs rex.install and reads the rhiza.lua it writes. The sized app's action has bounds = cell_bounds(ctx.block_id, 77, 21) and command { BIN, "tui", "app", "zz-sized" }, with no "launch" in it. The unsized app gets the 100x30 defaults, and no app action goes through tui launch. The rex-install.txtar script also passes (go test ./cmd/rhiza -run TestScripts/rex-install$) and greps bounds = cell_bounds(ctx.block_id, 100, 30). I read the diff: the Lua cell_bounds is a faithful port of rexuikit v0.2.0 surface.fitAxis/BoundsForCells, centred, with the Frame, Border and InverseMargin constants filled in from the Go package. It reads columns, rows, cell_width_px and cell_height_px, the field names of the terminal size result, and falls back to LAYER_BOUNDS. make fmt-check vet test passes. make self-check fails only on a stale .rhiza/rendered (actor claude-code-3b240678 and index.md), which comes from the parent run commit 6733c86. This leaf's commit touches no .rhiza file.

## Links

### Out

- verifies [rex-app-layer-sized](../capability/rhiza/rex-app-layer-sized.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
