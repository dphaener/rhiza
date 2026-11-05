# rex-install-path-plan

source · `01M6JBW2RZ1MW47W4GCFCMTNMP`

## quotes

```json
[]
```

## raw_sha256

9cf3e667ac3b1c474fa6494c3d5ad6e68b3a38110b907502d936c9de42d0ace6

## session

a5ad3bd3-86c2-4367-9548-01a1a4ba3593

## summary

Autonomous plan run for change set 01M6JBK. rhiza rex install wrote the installing host's absolute rhiza path into rhiza.lua, so app actions failed on a Rex pane hosted on the VM. One question was raised (how the actions find rhiza); the arbiter answered with Decision rex-actions-run-rhiza-from-path: run bare "rhiza" from the executing host's PATH, with an optional settings.bin override kept across reinstalls. One leaf planned under rex-install-command. No owner words: the run had no owner.

## title

Plan of rex install's rhiza path, autonomous run

## Links

### In

- derived-from [rex-actions-run-rhiza-from-path](../decision/rex-actions-run-rhiza-from-path.md)

## History

- `01M6JBKYEAHHCRYTMGBWTBZBBK` Fix rhiza rex install bakes the installing machine's absolute rhiza path into rhiza.lua (local BIN = "/Users/darinhaener/bin/rhiza"), but a Rex app runs on the host of its pane: with the Mac's Rex client working sessions on the VM's Rex server, the roadmap app failed with "rex-run: line 1: /Users/darinhaener/bin/rhiza: No such file or directory" (owner, 2026-11-03). The generated actions should run rhiza from the executing host's PATH (or resolve it there), not a path fixed at install time; a symlink at that path on the VM is the stopgap
