# native-actions-supersede-launcher

decision · `01M646RV3WM2456NBG2ZX59ZM2`

## reasoning

The owner chose to remove the picker, so the picker intent no longer holds; one new intent states what replaces it, and every leaf serves it.

## statement

The change set adds intent tui-native-actions under rhiza-tui, which supersedes tui-launcher (the picker); its four leaves serve tui-native-actions.

## wording

decided by the orchestrator acting as arbiter, following the owner's pick "Remove it" for the picker

## Links

### Out

- answers [native-intent-served](../question/native-intent-served.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
