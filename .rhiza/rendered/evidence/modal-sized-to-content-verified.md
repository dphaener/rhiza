# modal-sized-to-content-verified

evidence · `01M6N65M4AC0BR1CW5NK8F08SV` · actor [claude-code-911c09c8](../actor/claude-code-911c09c8.md)

## result

PASS

## statement

At 71044d36 ran go test -trimpath -count=1 ./internal/api -run TestRexInstall (ok), ./internal/tui/... (all ok) and ./cmd/rhiza -run TestScripts/rex-install (ok). Read the diff: launcher.App now carries a Component func and Size() returns that component's NaturalSize (the 100x30 default where an app has none). The new TestRexInstallOpensEveryAppAsACentredModalSizedToItsComponent (internal/api/rex_modal_test.go) registers apps hosting a Form, a List and a Dialog beside the real apps, installs rhiza.lua, runs every registered app's action under gopher-lua with a fake rex (240x70 window, 9x18 px cells), and checks each new_layer call opens in the focused block's window with bounds equal to surface.BoundsForCells(wc, NaturalSize, surface.Centered), smaller than the whole window, and that Size equals the component's NaturalSize; a grown List gets a larger modal on reinstall. This answers attempt 1's findings: sizes now come from NaturalSize, and the Lua's own bounds are compared against BoundsForCells. The rex-install testscript greps the written rhiza.lua under a fake rex for cell_bounds(ctx.block_id, 100, 30) and rejects any bounds line other than that call or LAYER_BOUNDS. Real apps do not host components yet, so they open at the default size.

## Links

### Out

- verifies [modal-sized-to-content](../capability/rhiza/modal-sized-to-content.md)

## History

- `01M6MYM5GB3Q8DBZST8V5AG6RQ` Modal shell: a shared package that opens a Rhiza app as a rexuikit surface Modal sized to its content, renders rexuikit ui components (Form, List, Dialog) with the rexuikit theme, and closes on Esc or when the app's job is done
