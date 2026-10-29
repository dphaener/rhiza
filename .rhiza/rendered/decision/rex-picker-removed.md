# rex-picker-removed

decision · `01M646QYNA203KTTF2VMRW575K`

## reasoning

The owner tried the picker and said it was not what they were looking for; with each app its own action and key, the picker adds a step and nothing else.

## statement

The all-apps picker (rhiza tui launch and the rhiza_launch action) is removed, as forge has none: each app opens only from its own key in the rhiza mode or from its own action in Rex's command palette. The capabilities that describe the picker are retired or revised to match.

## wording

the orchestrator's option, picked by the owner: "Remove it"

## Links

### Out

- answers [picker-kept](../question/picker-kept.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
