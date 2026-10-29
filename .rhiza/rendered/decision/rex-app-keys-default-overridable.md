# rex-app-keys-default-overridable

decision · `01M646QX1CVSV5JVY2A06362PW`

## reasoning

Defaults make the mode usable straight after install; a map in the settings block lets the user change any key, as the owner asked for the key that enters the mode.

## statement

rhiza ships a default mnemonic key for each app in the rhiza key mode, and the user can remap any app's key in rhiza.lua's settings block; a reinstall keeps the map.

## wording

the orchestrator's option, picked by the owner: "Defaults, overridable (Recommended)"

## Links

### Out

- answers [key-per-app](../question/key-per-app.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
