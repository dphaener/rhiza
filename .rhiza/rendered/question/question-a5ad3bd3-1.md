# question-a5ad3bd3-1

question · `01M6JBQQWRD69ZHZKC7Y7VAZPN`

## session

a5ad3bd3-86c2-4367-9548-01a1a4ba3593

## statement

rhiza rex install bakes the installing host's absolute rhiza path into rhiza.lua (local BIN = "/Users/.../bin/rhiza"), which breaks when the Rex pane runs on the VM. How should the generated app actions find rhiza?

## Links

### In

- answers [rex-actions-run-rhiza-from-path](../decision/rex-actions-run-rhiza-from-path.md)

## History

- `01M6JBKYEAHHCRYTMGBWTBZBBK` Fix rhiza rex install bakes the installing machine's absolute rhiza path into rhiza.lua (local BIN = "/Users/darinhaener/bin/rhiza"), but a Rex app runs on the host of its pane: with the Mac's Rex client working sessions on the VM's Rex server, the roadmap app failed with "rex-run: line 1: /Users/darinhaener/bin/rhiza: No such file or directory" (owner, 2026-11-03). The generated actions should run rhiza from the executing host's PATH (or resolve it there), not a path fixed at install time; a symlink at that path on the VM is the stopgap
