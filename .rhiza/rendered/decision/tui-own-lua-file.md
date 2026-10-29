# tui-own-lua-file

decision · `01M62VBWAKNS7FRMRZRNZRAYAM`

## reasoning

Done point tui-install-script requires that a second install leave the owner's own config untouched, and a separate file meets that without parsing init.lua. forge-cli already installs this way (forge.lua plus require("forge")), so the two tools share one pattern in the owner's Rex config.

## statement

The quick-launch script is its own file, rhiza.lua, beside Rex's init.lua, as forge keeps forge.lua. rhiza rex install writes it and adds one require("rhiza") line to init.lua if that line is missing. It changes no other line the owner wrote, and rhiza rex uninstall removes both.

## wording

the orchestrator's recommended option, taken under the owner's standing word: "i feel that I've made enough high level decisions" and "I'm just accepting the recommended answers anyway"

## Links

### Out

- answers [rex-config-file-choice](../question/rex-config-file-choice.md)

### In

- decided-by [rex-install-command](../capability/rhiza/rex-install-command.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
