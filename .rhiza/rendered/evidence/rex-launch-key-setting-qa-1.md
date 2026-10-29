# rex-launch-key-setting-qa-1

evidence · `01M648A25M01S6ECSZC7A35WE5` · actor [claude-code-fe283bac](../actor/claude-code-fe283bac.md)

## result

PASS

## statement

Same scratch install: the first install writes launch_key = "ctrl+a>a" and keys = {} in the settings block; I edited launch_key to ctrl+b>x and keys to { roadmap = "z" }, reran rhiza rex install --json: written [], unchanged [rhiza.lua, init.lua], file byte-identical (cmp), settings block kept.

## Links

### Out

- verifies [rex-launch-key-setting](../capability/rhiza/rex-launch-key-setting.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
