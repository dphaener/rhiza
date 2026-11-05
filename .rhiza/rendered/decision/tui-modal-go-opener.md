# tui-modal-go-opener

decision · `01M6N06QV22B773C7W5HHYP7KT`

## statement

Each app's modal is opened by rhiza itself: the Rex action runs rhiza, which opens the layer through rexuikit's surface over Rex's control connection, so the modal is sized to its content and can resize as the content changes

## wording

Go opener: the Rex action runs rhiza, which opens the modal through rexuikit's surface over Rex's control connection, so it can resize as its content changes, the agent's option, picked by the owner

## Links

### Out

- answers [modal-open-mechanism](../question/modal-open-mechanism.md)
- derived-from [tui-small-apps-interview](../source/tui-small-apps-interview.md)

### In

- decided-by [modal-sized-to-content](../capability/rhiza/modal-sized-to-content.md)

## History

- `01M6MYM5GB3Q8DBZST8V5AG6RQ` Modal shell: a shared package that opens a Rhiza app as a rexuikit surface Modal sized to its content, renders rexuikit ui components (Form, List, Dialog) with the rexuikit theme, and closes on Esc or when the app's job is done
