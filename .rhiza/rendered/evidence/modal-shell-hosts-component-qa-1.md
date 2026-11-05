# modal-shell-hosts-component-qa-1

evidence · `01M6N6BEJS6SX810TXSFW0D52Y` · actor [claude-code-074a669f](../actor/claude-code-074a669f.md)

## result

PASS

## statement

At c64eb6fa built go build -o $TMPDIR/rhiza-x ./cmd/rhiza and ran `rhiza-x tui app roadmap` and `rhiza-x tui app inbox` in a 120x40 tmux pane (no Rex): both now run through modalshell.Run (surface/program wrapper) and render their full views (roadmap milestones, inbox Overlaps/Failed runs), same as a base build at the merge-base. Scoped tests pass: go test -trimpath ./internal/tui/modalshell/... ./internal/tui/launcher/... and TestScripts/tui-launch-usage. No registered app hosts a rexuikit component yet, so component hosting itself is covered only by the package tests.

## Links

### Out

- verifies [modal-shell-hosts-component](../capability/rhiza/modal-shell-hosts-component.md)

## History

- `01M6MYM5GB3Q8DBZST8V5AG6RQ` Modal shell: a shared package that opens a Rhiza app as a rexuikit surface Modal sized to its content, renders rexuikit ui components (Form, List, Dialog) with the rexuikit theme, and closes on Esc or when the app's job is done
