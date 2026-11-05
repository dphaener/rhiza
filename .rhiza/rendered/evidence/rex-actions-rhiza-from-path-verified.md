# rex-actions-rhiza-from-path-verified

evidence · `01M6MEXF43MYKE6Y8CASBVVJ14` · actor [claude-code-0811a259](../actor/claude-code-0811a259.md)

## result

PASS

## statement

Verifier at 91fc3a3c: read the diff of internal/api/rex_install.go and rex_install_test.go. rexScript no longer takes a bin; rhiza.lua now writes `local BIN = "rhiza"` then `if type(settings.bin) == "string" and settings.bin ~= "" then BIN = settings.bin end` after the settings block, so a missing, empty or non-string bin leaves "rhiza" (Lua semantics; no interpreter on the host to execute it). The default settings block carries a commented-out bin entry. TestRexInstallScriptEntersTheKeyMode installs with ctx.Self="/opt/my tools/rhiza" and asserts that path appears nowhere in rhiza.lua; it then sets bin = "/home/me/go/bin/rhiza" and a new launch key in the settings block, reinstalls with ctx.Self="/elsewhere/rhiza", and asserts the file equals the edited bytes exactly and contains no "/elsewhere". TestRexInstallRoundTrip asserts a second install writes nothing and leaves rhiza.lua byte-identical. Ran `go test -trimpath ./internal/api -run TestRexInstall -count=1 -v` (all 9 PASS), `go test -trimpath ./cmd/rhiza -run 'TestScripts/(rex-install|daemon-install)$' -count=1` (ok), gofmt -l and go vet on internal/api (clean).

## Links

### Out

- verifies [rex-actions-rhiza-from-path](../capability/rhiza/rex-install-command/rex-actions-rhiza-from-path.md)

## History

- `01M6JBKYEAHHCRYTMGBWTBZBBK` Fix rhiza rex install bakes the installing machine's absolute rhiza path into rhiza.lua (local BIN = "/Users/darinhaener/bin/rhiza"), but a Rex app runs on the host of its pane: with the Mac's Rex client working sessions on the VM's Rex server, the roadmap app failed with "rex-run: line 1: /Users/darinhaener/bin/rhiza: No such file or directory" (owner, 2026-11-03). The generated actions should run rhiza from the executing host's PATH (or resolve it there), not a path fixed at install time; a symlink at that path on the VM is the stopgap
