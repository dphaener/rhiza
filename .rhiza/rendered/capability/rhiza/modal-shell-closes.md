# modal-shell-closes

capability · `01M6MZ92NEJNWBVABS0JB290GE` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in the shell's package drive it with bubbletea messages: Esc, FormCancelledMsg, ListCancelledMsg and the app's done message each return tea.Quit, any other key does not, and the result the app set is what the shell returns.

## statement

The modal shell closes the app on Esc, on the hosted component's cancel (ui.FormCancelledMsg, ui.ListCancelledMsg) and when the app reports its job done, and hands back the app's result for --once.

## Links

### Out

- decided-by [tui-small-modal](../../decision/tui-small-modal.md)

### In

- verifies [modal-shell-closes-verified](../../evidence/modal-shell-closes-verified.md)
- verifies [modal-shell-closes-qa-1](../../evidence/modal-shell-closes-qa-1.md)
- verifies [pr-290-modal-shell-for-rex-apps](../../evidence/pr-290-modal-shell-for-rex-apps.md)

## History

- `01M6MYM5GB3Q8DBZST8V5AG6RQ` Modal shell: a shared package that opens a Rhiza app as a rexuikit surface Modal sized to its content, renders rexuikit ui components (Form, List, Dialog) with the rexuikit theme, and closes on Esc or when the app's job is done
