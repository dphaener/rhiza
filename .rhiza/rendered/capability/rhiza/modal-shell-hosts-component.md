# modal-shell-hosts-component

capability · `01M6MZ8VSGET61TPSYJQ4AQCX1` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in the shell's package, with no Rex and no socket: a hosted Form, List and Dialog each render through rexuikit (the view carries the component's own render), a dark and a light tea.BackgroundColorMsg give theme.Dark and theme.Light tokens, and the program model is a surface/program.Model; the tui-launch-usage testscript still passes.

## statement

A shared modal shell package under internal/tui runs any Rhiza app as one rexuikit ui component (Form, List, Dialog or InfoPanel), themed from the terminal's background with rexuikit's ui/theme and wrapped in rexuikit's surface/program for focus, and rhiza tui app runs every app through it.

## Links

### Out

- decided-by [tui-small-modal](../../decision/tui-small-modal.md)
- decided-by [tui-rexuikit-components](../../decision/tui-rexuikit-components.md)

### In

- verifies [modal-shell-hosts-component-verified](../../evidence/modal-shell-hosts-component-verified.md)
- verifies [modal-shell-hosts-component-qa-1](../../evidence/modal-shell-hosts-component-qa-1.md)
- verifies [pr-290-modal-shell-for-rex-apps](../../evidence/pr-290-modal-shell-for-rex-apps.md)

## History

- `01M6MYM5GB3Q8DBZST8V5AG6RQ` Modal shell: a shared package that opens a Rhiza app as a rexuikit surface Modal sized to its content, renders rexuikit ui components (Form, List, Dialog) with the rexuikit theme, and closes on Esc or when the app's job is done
