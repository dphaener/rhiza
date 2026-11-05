# rex-actions-rhiza-from-path-qa-1

evidence · `01M6MF2RQWRMQYXPC67G64Y44Q` · actor [claude-code-8c45d7fc](../actor/claude-code-8c45d7fc.md)

## result

PASS

## statement

QA at a6c27bb4 with the built binary (go build -trimpath ./cmd/rhiza), REX_CONFIG pointing at a scratch folder and a fake rex on PATH for the reload: `rhiza rex install` wrote rhiza.lua with `local BIN = "rhiza"` plus the settings.bin guard and a commented `-- bin` example in the settings block; the binary's own scratch path appears nowhere in the file. A second install reported both files unchanged and cmp showed rhiza.lua byte-identical. With bin = "/home/me/go/bin/rhiza" and launch_key "ctrl+a g" set in the settings block, a reinstall left the file byte-identical (cmp). An old-style rhiza.lua with `local BIN = "/Users/darinhaener/bin/rhiza"` was rewritten on reinstall to the PATH form, keeping the owner's settings. The real Rex (`rex config check` with REX_CONFIG on the scratch init.lua) loads the generated file (24 bindings, 1 mode, 22 actions). A probe copy raising error("PROBE BIN="..BIN) after the guard, loaded by real Rex, gave BIN=rhiza for a missing bin, "", 42 and false, and BIN=/home/me/go/bin/rhiza when that string is set. go test ./internal/api -run TestRexInstall and TestScripts/rex-install pass.

## Links

### Out

- verifies [rex-actions-rhiza-from-path](../capability/rhiza/rex-install-command/rex-actions-rhiza-from-path.md)

## History

- `01M6JBKYEAHHCRYTMGBWTBZBBK` Fix rhiza rex install bakes the installing machine's absolute rhiza path into rhiza.lua (local BIN = "/Users/darinhaener/bin/rhiza"), but a Rex app runs on the host of its pane: with the Mac's Rex client working sessions on the VM's Rex server, the roadmap app failed with "rex-run: line 1: /Users/darinhaener/bin/rhiza: No such file or directory" (owner, 2026-11-03). The generated actions should run rhiza from the executing host's PATH (or resolve it there), not a path fixed at install time; a symlink at that path on the VM is the stopgap
