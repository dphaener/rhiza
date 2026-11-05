# tui-install-script

intent · `01M62JMHXNM6W4XBGP85W06YY8` · parent [rhiza-tui](../rhiza-tui.md)

## statement

rhiza rex install writes Rhiza's app actions and key mode as a Lua script into the Rex config and reloads Rex, and a second run leaves the owner's own config untouched

## Links

### In

- serves [rex-install-command](../../capability/rhiza/rex-install-command.md)
- serves change set `01M6JBKYEAHHCRYTMGBWTBZBBK` (Fix rhiza rex install bakes the installing machine's absolute rhiza path into rhiza.lua (local BIN = "/Users/darinhaener/bin/rhiza"), but a Rex app runs on the host of its pane: with the Mac's Rex client working sessions on the VM's Rex server, the roadmap app failed with "rex-run: line 1: /Users/darinhaener/bin/rhiza: No such file or directory" (owner, 2026-11-03). The generated actions should run rhiza from the executing host's PATH (or resolve it there), not a path fixed at install time; a symlink at that path on the VM is the stopgap)

## History

- `01M62J0ED188JTTR7R9709Y989` Rhiza has a terminal UI built on Rex, with rexuikit, in which common rhiza commands run as tiny apps in modals or pop-ups, alongside a view of the model.
- `01M6489HCR688XKH99YV9E2VGY` Remove the all-apps picker, as decision rex-picker-removed says: rhiza tui launch and the rhiza_launch action written by rhiza rex install go, and the capabilities that describe the picker are retired or revised, so each app opens only from its own key in the rhiza key mode or its own action in Rex's command palette.
