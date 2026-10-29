# rex-install-command

capability · `01M62JP4KDNBZKH1YT5NZRB347` · parent [rhiza](../rhiza.md)

## acceptance

With a fake Rex config directory and a fake rex, install writes rhiza.lua and the require line. A second install leaves both files byte-identical. Every line the owner wrote in init.lua is kept byte for byte. uninstall restores init.lua as it was, and no test touches real Rex.

## statement

rhiza rex install writes rhiza.lua beside Rex's init.lua, adds require("rhiza") to init.lua if it is missing, and reloads Rex. A second run changes nothing the owner wrote, keeps the settings in rhiza.lua, and rhiza rex uninstall removes both.

## Links

### Out

- decided-by [tui-many-mini-apps](../../decision/tui-many-mini-apps.md)
- decided-by [tui-own-lua-file](../../decision/tui-own-lua-file.md)
- implements [tui-install-script](../../intent/rhiza-tui/tui-install-script.md)

### In

- blocks [rex-config-file-choice](../../question/rex-config-file-choice.md)
- verifies [rex-install-command-verified](../../evidence/rex-install-command-verified.md)
- verifies [rex-install-command-qa-1](../../evidence/rex-install-command-qa-1.md)
- verifies [rex-install-command-qa-2](../../evidence/rex-install-command-qa-2.md)
- verifies [rex-install-command-qa-3](../../evidence/rex-install-command-qa-3.md)
- verifies [pr-65-tui-shell-and-read-apps](../../evidence/pr-65-tui-shell-and-read-apps.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
