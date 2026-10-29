# rex-app-own-action

decision · `01M644QVPPBD4ZKTT1YP9SMGAS`

## quote

the forge-cli opens the forge quick actions entirely differently, with a lua script, natively in Rex.

## statement

Each rhiza app is its own Rex action, opened from the command palette and by the key mode, with no picker in between

## Links

### Out

- derived-from [native-rex-launch-owner-words](../source/native-rex-launch-owner-words.md)

### In

- decided-by [rex-app-actions](../capability/rhiza/rex-app-actions.md)
- decided-by [rex-app-layer-sized](../capability/rhiza/rex-app-layer-sized.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
