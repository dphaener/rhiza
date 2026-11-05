# modal-open-mechanism

question · `01M6MZ9FFYJF8W6MDPVS2T2SVK`

## statement

How does an app's Rex action open its Modal, and does the Modal follow its content's size as it changes (an inbox list growing or shrinking)? Either the Lua action keeps opening the layer with rex.session.new_layer, sized from cells fixed when rhiza rex install writes it, or a Go opener calls rexuikit surface.Open (Kind Modal, Cells from NaturalSize) over Rex's control connection and resizes with Handle.Resize.

## Links

### Out

- blocks [modal-sized-to-content](../capability/rhiza/modal-sized-to-content.md)

### In

- answers [tui-modal-go-opener](../decision/tui-modal-go-opener.md)

## History

- `01M6MYM5GB3Q8DBZST8V5AG6RQ` Modal shell: a shared package that opens a Rhiza app as a rexuikit surface Modal sized to its content, renders rexuikit ui components (Form, List, Dialog) with the rexuikit theme, and closes on Esc or when the app's job is done
