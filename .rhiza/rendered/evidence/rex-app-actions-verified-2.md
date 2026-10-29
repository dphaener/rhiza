# rex-app-actions-verified-2

evidence · `01M648YY1ST16ZB9MAE7GE2419` · actor [claude-code-c1be1bd3](../actor/claude-code-c1be1bd3.md)

## result

PASS

## statement

Verifier, attempt 2 at b0bf33a: go test -count=1 ./cmd/rhiza -run TestScripts/rex-install passes (rex-install.txtar checks per-app actions titled with the app name running tui app <name>, now incl. name = rhiza_app_session_ls and no action name outside [a-z0-9_]); go test ./internal/api -run TestRexInstall passes, incl. TestRexInstallRegistersAnActionPerApp (one action per readapps/writeapps/modelapp app; a registered zz-extra app found after reinstall as rhiza_app_zz_extra, title zz-extra, command tui app zz-extra). go test -count=1 ./... all ok. Built the binary, ran rhiza rex install into a scratch XDG_CONFIG_HOME, and real `rex config check` on its init.lua exits 0: 17 actions (16 apps plus rhiza_launch), 18 bindings, 1 mode, no diagnostics; the 7 hyphenated apps (change-show, record-*, run-status, session-ls) are now defined with '_' names, and the key mode's rex.bind maps names the same way via gsub. make check's self-check fails only on stale .rhiza/rendered files from run writes, not on code.

## Links

### Out

- verifies [rex-app-actions](../capability/rhiza/rex-app-actions.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
