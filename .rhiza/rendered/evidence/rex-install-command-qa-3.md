# rex-install-command-qa-3

evidence · `01M631J6B9ABWKNZG1V5EEWDY7` · actor [claude-code-2ac42aad](../actor/claude-code-2ac42aad.md)

## result

PASS

## statement

QA at d02d9ce (cycle 3) with a built rhiza, REX_CONFIG at a scratch init.lua (owner lines incl. trailing spaces) and a fake rex first on PATH. rex install wrote rhiza.lua (settings launch_key ctrl+a a, BIN, block_cwd helper, rhiza_launch action passing cwd, bind) and appended one require("rhiza") -- added by rhiza rex install line; the fake rex logged config reload. A second install reported both unchanged and cmp found both files byte-identical. uninstall removed both and cmp shows init.lua equal to the owner's original. A real scratch Rex server loading the installed config bound ctrl+a a to rhiza_launch.

## Links

### Out

- verifies [rex-install-command](../capability/rhiza/rex-install-command.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
