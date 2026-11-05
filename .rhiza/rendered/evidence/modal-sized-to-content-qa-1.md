# modal-sized-to-content-qa-1

evidence · `01M6N6BSWGXXCC9V3B8HDEJ6XK` · actor [claude-code-074a669f](../actor/claude-code-074a669f.md)

## result

PASS

## statement

At c64eb6fa ran the built `rhiza-x rex install` with REX_CONFIG in a scratch folder and a fake rex on PATH: it wrote rhiza.lua and init.lua and the fake rex logged `config reload`. Every action's layer is `bounds = cell_bounds(ctx.block_id, 100, 30)` (the default, since no app registers a Component yet), and cell_bounds centres a span capped at 1 from the window's cells, so no action opens full-screen. go test -trimpath ./internal/api -run TestRexInstall (incl. the Form/List/Dialog NaturalSize modal test) and TestScripts/rex-install pass.

## Links

### Out

- verifies [modal-sized-to-content](../capability/rhiza/modal-sized-to-content.md)

## History

- `01M6MYM5GB3Q8DBZST8V5AG6RQ` Modal shell: a shared package that opens a Rhiza app as a rexuikit surface Modal sized to its content, renders rexuikit ui components (Form, List, Dialog) with the rexuikit theme, and closes on Esc or when the app's job is done
