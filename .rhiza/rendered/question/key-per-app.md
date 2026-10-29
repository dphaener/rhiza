# key-per-app

question · `01M644QYKDPX7RCMG3MVFFD461`

## statement

Which key does each app get in the rhiza key mode: its first letter (several apps share one) or a key the owner sets per app in the settings block?

## Links

### Out

- blocks [rex-key-mode](../capability/rhiza/rex-key-mode.md)

### In

- answers [rex-app-keys-default-overridable](../decision/rex-app-keys-default-overridable.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
