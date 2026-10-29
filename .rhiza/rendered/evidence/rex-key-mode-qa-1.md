# rex-key-mode-qa-1

evidence · `01M648A0XW09CX6C0YNCSF35SV` · actor [claude-code-fe283bac](../actor/claude-code-fe283bac.md)

## result

FAIL

## statement

Same install, real rex config check: the launch key binds client.mode.enter {name=rhiza, once=true}, mode rhiza is exclusive with escape to exit, and DEFAULT_KEYS gives each of the 16 apps a distinct key (accept=a ... why=w), 18 bindings, 1 mode. But 7 of those keys bind to rhiza_app_<hyphenated-name> actions that Rex refused to define, so their keys open nothing. With hyphens mapped to underscores in both the names and the bind, config check reports 17 actions, 18 bindings, no diagnostics.

## Links

### Out

- verifies [rex-key-mode](../capability/rhiza/rex-key-mode.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
