# modal-sized-to-content

capability · `01M6MZ97WECAS9XNZNBF3T2JV3` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in internal/tui/modalopen drive the opener against a fake Rex: it opens the app with rexuikit surface.Open as a Modal over the invoking block's window, with new_layer bounds equal to surface.BoundsForCells of the content's NaturalSize centred (zero Anchor) and never the whole window, and each new size the modal shell reports resizes the layer with Handle.Resize to BoundsForCells of it; a test in internal/api runs every registered app's Rex action from the written rhiza.lua under a fake rex and checks it opens no layer itself but starts rhiza tui open <app> detached, and the rex-install testscript greps that from the written rhiza.lua with a fake rex.

## statement

Each app opens as a centred rexuikit Modal over the current pane, sized to its component's NaturalSize and never full-screen.

## Links

### Out

- decided-by [tui-small-modal](../../decision/tui-small-modal.md)
- decided-by [tui-rexuikit-components](../../decision/tui-rexuikit-components.md)
- decided-by [tui-modal-go-opener](../../decision/tui-modal-go-opener.md)

### In

- blocks [modal-open-mechanism](../../question/modal-open-mechanism.md)
- verifies [modal-sized-to-content-verified](../../evidence/modal-sized-to-content-verified.md)
- verifies [modal-sized-to-content-qa-1](../../evidence/modal-sized-to-content-qa-1.md)
- verifies [modal-sized-to-content-verified-2](../../evidence/modal-sized-to-content-verified-2.md)
- verifies [pr-290-modal-shell-for-rex-apps](../../evidence/pr-290-modal-shell-for-rex-apps.md)

## History

- `01M6MYM5GB3Q8DBZST8V5AG6RQ` Modal shell: a shared package that opens a Rhiza app as a rexuikit surface Modal sized to its content, renders rexuikit ui components (Form, List, Dialog) with the rexuikit theme, and closes on Esc or when the app's job is done
