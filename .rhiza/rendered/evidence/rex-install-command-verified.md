# rex-install-command-verified

evidence · `01M62WRG0CVY6M0H6ZMCE97ZED` · actor [claude-code-9ffb274d](../actor/claude-code-9ffb274d.md)

## result

PASS

## statement

At fe0f742 I ran go test -count=1 ./internal/api -run TestRexInstall (4 tests PASS) and the testscript TestScripts/rex-install plus describe-golden (PASS), then make check with RHIZA_ROOT/CHANGE/RUN unset (exit 0). The rex-install script sets REX_CONFIG to $WORK/conf/init.lua and puts a fake rex first on PATH; the unit tests inject Getenv and Exec, so no test touches real Rex. First install writes rhiza.lua (marker line, launch_key ctrl+a a) and exactly one require("rhiza") line, and keeps the owner's lines; rex.log records one 'config reload'. A second install reports written [] and cmp shows both files byte-identical. An edited setting survives a reinstall. Uninstall restores init.lua byte for byte (cmp against owner.lua; the unit test also covers no trailing newline, an empty file and odd spacing) and removes rhiza.lua. In the diff, rexAddRequire only appends and rexRemoveRequire removes only the tagged line; an owner's own require is left alone, and a rhiza.lua install did not write is a conflict unless --force.

## Links

### Out

- verifies [rex-install-command](../capability/rhiza/rex-install-command.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
