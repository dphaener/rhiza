# rex-actions-run-rhiza-from-path

decision · `01M6JBT72RR7K9DS2FDF6J2DAG`

## reasoning

The change set's intent, quoting the owner's failure on 2026-11-03, already says the actions should run rhiza from the executing host's PATH (or resolve it there), not a path fixed at install time, and calls the VM symlink only a stopgap. A Rex app runs on its pane's host, so no path the installer knows is right in general; only the executing host can resolve it. The failure text (rex-run: line 1: ...) shows the command goes through a shell on that host, so a bare name is resolved from its PATH. Resolving in Lua at runtime (probing paths) adds code with no gain over PATH lookup. Because a pane's environment may lack ~/bin on PATH, an escape hatch is needed; the settings block (rexSettingsBegin/End in internal/api/rex_install.go) is already the part a reinstall preserves, so a bin setting there is the smallest fit and keeps per-host overrides across make update. Tests in internal/api/rex_install_test.go that assert local BIN = <ctx.Self> change to assert the PATH default and a kept bin override.

## statement

rhiza rex install no longer writes the installing host's absolute path into rhiza.lua. The generated app actions run the bare command "rhiza", resolved on the PATH of the host that runs the pane, unless the owner sets an optional bin entry in the settings block (e.g. bin = "/home/me/bin/rhiza"), which a reinstall keeps like launch_key and keys. So the script reads local BIN = (type(settings.bin) == "string" and settings.bin ~= "") and settings.bin or "rhiza" (a missing, empty or non-string bin falls back to "rhiza", the same guard launch_key uses), and ctx.Self no longer feeds rhiza.lua.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-a5ad3bd3-1](../question/question-a5ad3bd3-1.md)
- derived-from [rex-install-path-plan](../source/rex-install-path-plan.md)

### In

- decided-by [rex-actions-rhiza-from-path](../capability/rhiza/rex-install-command/rex-actions-rhiza-from-path.md)

## History

- `01M6JBKYEAHHCRYTMGBWTBZBBK` Fix rhiza rex install bakes the installing machine's absolute rhiza path into rhiza.lua (local BIN = "/Users/darinhaener/bin/rhiza"), but a Rex app runs on the host of its pane: with the Mac's Rex client working sessions on the VM's Rex server, the roadmap app failed with "rex-run: line 1: /Users/darinhaener/bin/rhiza: No such file or directory" (owner, 2026-11-03). The generated actions should run rhiza from the executing host's PATH (or resolve it there), not a path fixed at install time; a symlink at that path on the VM is the stopgap
