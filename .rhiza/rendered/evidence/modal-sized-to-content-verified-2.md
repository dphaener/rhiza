# modal-sized-to-content-verified-2

evidence · `01M6N7NST4T1S2Q4G7X18KTBRP` · actor [claude-code-2a05f8b0](../actor/claude-code-2a05f8b0.md)

## result

PASS

## statement

Verifier, attempt 2, at 73831be9: read the diff and ran go test -race ./internal/tui/modalopen/... ./internal/tui/modalshell/... (ok), go test ./internal/api/... (ok), and TestScripts/(rex-install|tui-open-usage|tui-read-apps|describe-golden) (ok); go vet and gofmt clean. TestOpenOpensACentredModalAndFollowsItsContent drives modalopen.Open against a fake Rex. It sees one new_layer in the opener block's window w1, not the active w2, with Focus set and bounds equal to surface.BoundsForCells(window, 60, 12, Centered) (<1 on both axes, so never the whole window). Later sizes 90x30 and 40x8 each produce one Handle.Resize to BoundsForCells of that size, and a repeat of an unchanged size produces none. Open returns on block_closed. modalshell's size reporter writes NaturalSize on Init and on change, and WatchSize relays it. TestRexInstallActionsStartTheModalOpener runs every registered app's action from the written rhiza.lua under a fake rex. Each opens no layer and calls rex.exec(rhiza tui open <app> --session s --block b) with detach and the block's cwd. rex-install.txtar greps rex.exec/open_argv and that new_layer, bounds and tui app are absent. Not checked: real Rex (rex.exec's existence is assumed), and no registered app sets Component yet, so real apps open at the 100x30 default.

## Links

### Out

- verifies [modal-sized-to-content](../capability/rhiza/modal-sized-to-content.md)

## History

- `01M6MYM5GB3Q8DBZST8V5AG6RQ` Modal shell: a shared package that opens a Rhiza app as a rexuikit surface Modal sized to its content, renders rexuikit ui components (Form, List, Dialog) with the rexuikit theme, and closes on Esc or when the app's job is done
