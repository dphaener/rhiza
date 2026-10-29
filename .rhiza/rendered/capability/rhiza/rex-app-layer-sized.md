# rex-app-layer-sized

capability · `01M644QWTP8H6S0B5RS70QCBM2` · parent [rhiza](../rhiza.md)

## acceptance

A rex-install script test: for an app registered with Cols and Rows, its action's layer bounds come from those cells, and its action runs rhiza tui app <name> with no tui launch step

## statement

Each app's action opens it straight away in a Rex layer sized to that app's own columns and rows, with no list step first

## Links

### Out

- serves [tui-native-actions](../../intent/rhiza-tui/tui-native-actions.md)
- decided-by [rex-app-own-action](../../decision/rex-app-own-action.md)

### In

- verifies [rex-app-layer-sized-verified](../../evidence/rex-app-layer-sized-verified.md)
- verifies [rex-app-layer-sized-qa-1](../../evidence/rex-app-layer-sized-qa-1.md)
- verifies [rex-app-layer-sized-verified-2](../../evidence/rex-app-layer-sized-verified-2.md)
- verifies [rex-app-layer-sized-qa-2](../../evidence/rex-app-layer-sized-qa-2.md)
- verifies [pr-77-native-rex-actions-for-apps](../../evidence/pr-77-native-rex-actions-for-apps.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
