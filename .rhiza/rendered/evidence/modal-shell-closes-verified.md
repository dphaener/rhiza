# modal-shell-closes-verified

evidence · `01M6N570ENYFC6CKZKMSRCN2NM` · actor [claude-code-b741bd1b](../actor/claude-code-b741bd1b.md)

## result

PASS

## statement

At e7bb73fc ran go test -trimpath -count=1 -v ./internal/tui/modalshell/...: all pass. TestCloses drives a hosted Dialog shell with bubbletea messages and asserts Esc (tea.KeyPressMsg KeyEscape), ui.FormCancelledMsg, ui.ListCancelledMsg and DoneMsg each return a cmd yielding tea.QuitMsg and set Closed, while a plain 'a' key and a WindowSizeMsg do not. TestClosesEveryHost repeats the plain-key/Esc check across Form, List, Dialog and InfoPanel hosts. TestResult checks DoneMsg{Result:"out"} makes Result() return "out", DoneMsg{Err} returns that error, and a cancel keeps no result. Read the diff and Shell.Update in internal/tui/modalshell/modalshell.go: it quits on DoneMsg (keeping Result/Err), the cancel messages and esc, and forwards everything else to the component; Run returns s.Result(). go vet and gofmt clean.

## Links

### Out

- verifies [modal-shell-closes](../capability/rhiza/modal-shell-closes.md)

## History

- `01M6MYM5GB3Q8DBZST8V5AG6RQ` Modal shell: a shared package that opens a Rhiza app as a rexuikit surface Modal sized to its content, renders rexuikit ui components (Form, List, Dialog) with the rexuikit theme, and closes on Esc or when the app's job is done
