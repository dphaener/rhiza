# native-rex-prefix-configurable

decision · `01M644AP63VKTRFH6ZMW53XYGM`

## reasoning

The owner wants users to map the key to whatever they like, with ctrl+a>a as the default. A setting in rhiza.lua is how the current launch_key already works and how forge's prefix works.

## statement

The key that enters the rhiza key mode is a setting the user can map to any key, kept in the settings block of rhiza.lua so a reinstall keeps it; its default is ctrl+a>a.

## wording

The owner: "I think that key should be configurable, so users can map it to whatever key that want, but ctrl-a>a is great as a default."

## Links

### Out

- derived-from [native-rex-key-owner-words](../source/native-rex-key-owner-words.md)

### In

- decided-by change set `01M643X7CEAQP16925KTAM06GT` (Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.)
- decided-by [rex-key-mode](../capability/rhiza/rex-key-mode.md)
- decided-by [rex-launch-key-setting](../capability/rhiza/rex-launch-key-setting.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
