# rex-key-mode

capability · `01M644QXDASM0DPKA3VKKFGJD7` · parent [rhiza](../rhiza.md)

## acceptance

A rex-install script test: rhiza.lua binds one key in the key mode for each registered app, each bound to that app's action, and the keys are the ones the owner chose (see the key question)

## statement

Pressing the launch key enters a rhiza key mode with one key per app, and each key opens that app's action

## Links

### Out

- serves [tui-native-actions](../../intent/rhiza-tui/tui-native-actions.md)
- decided-by [native-rex-prefix-configurable](../../decision/native-rex-prefix-configurable.md)

### In

- blocks [key-per-app](../../question/key-per-app.md)
- verifies [rex-key-mode-verified](../../evidence/rex-key-mode-verified.md)
- verifies [rex-key-mode-qa-1](../../evidence/rex-key-mode-qa-1.md)
- verifies [rex-key-mode-verified-2](../../evidence/rex-key-mode-verified-2.md)
- verifies [rex-key-mode-qa-2](../../evidence/rex-key-mode-qa-2.md)
- verifies [pr-77-native-rex-actions-for-apps](../../evidence/pr-77-native-rex-actions-for-apps.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
