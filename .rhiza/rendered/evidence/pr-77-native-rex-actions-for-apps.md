# pr-77-native-rex-actions-for-apps

evidence · `01M64AZXR8841W9JG68H3W5SA8` · actor [darin-haener](../actor/darin-haener.md)

## head

a33d3385da54ad383212dcd2b4fb135745d9d7e1

## pr

https://github.com/dphaener/rhiza/pull/77

## result

1 check passed: check SUCCESS

## statement

Pull request #77 (https://github.com/dphaener/rhiza/pull/77) passed CI at a33d338

## status

verified

## Links

### Out

- verifies [rex-app-actions](../capability/rhiza/rex-app-actions.md)
- verifies [rex-app-layer-sized](../capability/rhiza/rex-app-layer-sized.md)
- verifies [rex-key-mode](../capability/rhiza/rex-key-mode.md)
- verifies [rex-launch-key-setting](../capability/rhiza/rex-launch-key-setting.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
