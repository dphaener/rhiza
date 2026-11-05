# pr-266-fix-rhiza-rex-install-bakes

evidence · `01M6MGAB4C3DC857AMPWJC097D` · actor [dphaener](../actor/dphaener.md)

## head

9096eb26248d13d96d3eca7b0f93e5b39adc1add

## pr

https://github.com/dphaener/rhiza/pull/266

## result

1 check passed: check SUCCESS

## statement

Pull request #266 (https://github.com/dphaener/rhiza/pull/266) passed CI at 9096eb2

## status

verified

## Links

### Out

- verifies [rex-actions-rhiza-from-path](../capability/rhiza/rex-install-command/rex-actions-rhiza-from-path.md)

## History

- `01M6JBKYEAHHCRYTMGBWTBZBBK` Fix rhiza rex install bakes the installing machine's absolute rhiza path into rhiza.lua (local BIN = "/Users/darinhaener/bin/rhiza"), but a Rex app runs on the host of its pane: with the Mac's Rex client working sessions on the VM's Rex server, the roadmap app failed with "rex-run: line 1: /Users/darinhaener/bin/rhiza: No such file or directory" (owner, 2026-11-03). The generated actions should run rhiza from the executing host's PATH (or resolve it there), not a path fixed at install time; a symlink at that path on the VM is the stopgap
