# rex-install-command-qa-1

evidence · `01M62YXZKKXCD34HTF84ST7NZW` · actor [claude-code-655e5878](../actor/claude-code-655e5878.md)

## result

PASS

## statement

QA at 4a03cf7 with a built rhiza (go build to a scratch dir), REX_CONFIG pointed at a scratch init.lua and a fake rex first on PATH. First rex install wrote rhiza.lua (marker, launch_key ctrl+a a, BIN = the binary) and appended one require("rhiza") -- added by rhiza rex install line, keeping the owner's lines incl. trailing spaces; fake rex logged config reload. A second install reported written [] and cmp showed both files byte-identical. An edited launch_key survived a reinstall. uninstall removed rhiza.lua and cmp showed init.lua equal to the owner's original. With no init.lua, install created it and uninstall removed it again; a second uninstall says nothing to remove. A foreign rhiza.lua is refused (exit 3, next: --force). The written init.lua plus rhiza.lua loaded in luajit against a fake rex: bind ctrl+a>a rhiza_launch, whose run calls new_layer with flavor com.superlogical.terminal.shell running <BIN> tui launch, focus true. Real Rex config untouched.

## Links

### Out

- verifies [rex-install-command](../capability/rhiza/rex-install-command.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
