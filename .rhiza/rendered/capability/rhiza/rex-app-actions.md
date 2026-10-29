# rex-app-actions

capability · `01M644QW8GJBD8KTJ8Q00MMNYY` · parent [rhiza](../rhiza.md)

## acceptance

A rex-install script test: after rhiza rex install, rhiza.lua has one rex.action per registered app (readapps, writeapps and modelapp), each titled with the app's name and running rhiza tui app <name>; a test that registers one more app and reinstalls finds its action in rhiza.lua

## statement

Each rhiza app is a Rex action of its own in the command palette, titled with the app's name, and running it runs that app

## Links

### Out

- serves [tui-native-actions](../../intent/rhiza-tui/tui-native-actions.md)
- decided-by [rex-app-own-action](../../decision/rex-app-own-action.md)

### In

- verifies [rex-app-actions-verified](../../evidence/rex-app-actions-verified.md)
- verifies [rex-app-actions-qa-1](../../evidence/rex-app-actions-qa-1.md)
- verifies [rex-app-actions-verified-2](../../evidence/rex-app-actions-verified-2.md)
- verifies [rex-app-actions-qa-2](../../evidence/rex-app-actions-qa-2.md)
- verifies [pr-77-native-rex-actions-for-apps](../../evidence/pr-77-native-rex-actions-for-apps.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
