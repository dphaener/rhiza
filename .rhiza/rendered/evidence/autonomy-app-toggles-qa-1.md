# autonomy-app-toggles-qa-1

evidence · `01M6NB1H7C8KN0V9PMH17P1D5S` · actor [claude-code-fb70f168](../actor/claude-code-fb70f168.md)

## result

PASS

## statement

QA at 1de07e73 with the built binary (go build -trimpath ./cmd/rhiza) in a scratch git repo with rhiza init, driving the real TUI in tmux: rhiza tui app autonomy drew one form with Plan/Accept/Ship/Sign-off, each auto/owner, prefilled from daemon settings (plan owner default, accept auto from git config, ship auto, signoff auto). Changing plan to auto and ship to owner and pressing enter through to the last field closed the app with exit 0; .git/config then held rhiza.autonomy.plan=auto and rhiza.autonomy.ship=owner, accept stayed as before and signoff was never written; rhiza daemon settings --json showed both as git config. Reopening prefilled the new values; changing plan then Esc exited 0 with .git/config byte-identical. rhiza tui app autonomy --once cmp-matched rhiza daemon settings --json. go test -trimpath ./internal/tui/autonomyapp/... and TestScripts/tui-autonomy-app passed.

## Links

### Out

- verifies [autonomy-app-toggles](../capability/rhiza/autonomy-app-toggles.md)

## History

- `01M6MYPA335J1S3952RJAGZJN9` Autonomy app: the plan, accept, ship and sign-off toggles (rhiza.autonomy.*), rebuilt on the modal shell over daemon.settings and daemon.set, replacing internal/tui/autonomyapp and retiring the capabilities of the old autonomy app it replaces
