# tui-native-actions

intent · `01M646RP92V3GJ99FNAPJCEKZN` · parent [rhiza-tui](../rhiza-tui.md)

## statement

Each rhiza app opens straight away from its own key in the rhiza key mode (ctrl+a>a then the app's key, both remappable in rhiza.lua) or from its own action in Rex's command palette, in a layer sized to it, with no picker in between

## Links

### Out

- supersedes [tui-launcher](tui-launcher.md)

### In

- serves [rex-app-actions](../../capability/rhiza/rex-app-actions.md)
- serves [rex-app-layer-sized](../../capability/rhiza/rex-app-layer-sized.md)
- serves [rex-key-mode](../../capability/rhiza/rex-key-mode.md)
- serves [rex-launch-key-setting](../../capability/rhiza/rex-launch-key-setting.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
