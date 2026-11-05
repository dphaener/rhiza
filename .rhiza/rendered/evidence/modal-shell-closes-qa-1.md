# modal-shell-closes-qa-1

evidence · `01M6N6BM8MK4N3EB82N3CZ2234` · actor [claude-code-074a669f](../actor/claude-code-074a669f.md)

## result

PASS

## statement

At c64eb6fa, with the built binary in a tmux pane, sent Escape to `rhiza-x tui app roadmap` and to `rhiza-x tui app inbox`: each exited with EXIT=0 within 2s, matching a base build at the merge-base. go test -trimpath -count=1 ./internal/tui/modalshell/... passes (Esc, FormCancelledMsg, ListCancelledMsg and DoneMsg quit; a plain key does not; Result returns the app's result).

## Links

### Out

- verifies [modal-shell-closes](../capability/rhiza/modal-shell-closes.md)

## History

- `01M6MYM5GB3Q8DBZST8V5AG6RQ` Modal shell: a shared package that opens a Rhiza app as a rexuikit surface Modal sized to its content, renders rexuikit ui components (Form, List, Dialog) with the rexuikit theme, and closes on Esc or when the app's job is done
