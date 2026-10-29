# rex-install-command-qa-2

evidence · `01M630C562VF8PBTRY6FBSN0VZ` · actor [claude-code-f3ca46cd](../actor/claude-code-f3ca46cd.md)

## result

PASS

## statement

QA at c7a97c9 (cycle 2) with a built rhiza, REX_CONFIG at a scratch init.lua (owner lines incl. trailing spaces) and a fake rex first on PATH. rex install wrote rhiza.lua (settings launch_key ctrl+a a, BIN = the built binary, rhiza_launch action, bind) and appended one require("rhiza") -- added by rhiza rex install line; the fake rex logged `config reload`. A second install reported both files unchanged and cmp found them byte-identical. rex uninstall removed both, and cmp shows init.lua equal to the original. Real Rex loaded the written rhiza.lua via rex do and its action opened the layer.

## Links

### Out

- verifies [rex-install-command](../capability/rhiza/rex-install-command.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
