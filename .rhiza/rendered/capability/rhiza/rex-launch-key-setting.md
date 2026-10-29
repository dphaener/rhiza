# rex-launch-key-setting

capability · `01M644QXZW01JJZ0SJRZA71GWV` · parent [rhiza](../rhiza.md)

## acceptance

A rex-install script test: a first install writes launch_key = "ctrl+a>a" (today the code writes ctrl+a a, the same key in the setting's spacing, so the test and the default change together); the owner then edits launch_key, a reinstall leaves that edited line byte for byte, and the settings block is kept

## statement

The key that enters the rhiza key mode is launch_key in rhiza.lua's settings block, default ctrl+a>a, and a reinstall keeps the value the owner set

## Links

### Out

- serves [tui-native-actions](../../intent/rhiza-tui/tui-native-actions.md)
- decided-by [native-rex-prefix-configurable](../../decision/native-rex-prefix-configurable.md)

### In

- verifies [rex-launch-key-setting-verified](../../evidence/rex-launch-key-setting-verified.md)
- verifies [rex-launch-key-setting-qa-1](../../evidence/rex-launch-key-setting-qa-1.md)
- verifies [rex-launch-key-setting-verified-2](../../evidence/rex-launch-key-setting-verified-2.md)
- verifies [rex-launch-key-setting-qa-2](../../evidence/rex-launch-key-setting-qa-2.md)
- verifies [pr-77-native-rex-actions-for-apps](../../evidence/pr-77-native-rex-actions-for-apps.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
