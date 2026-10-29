# native-supersedes-tui-launcher

decision · `01M6459DKBTY66HJ8PXPWVMHQH`

## reasoning

This change set's intent says each app opens straight away 'with no picker in between', and its leaves (rex-key-mode, rex-launch-key-setting) put the key mode on the launch key. The owner's own words in native-rex-prefix-configurable ask for ctrl-a>a as the default for the key that enters that mode, so the key and the no-picker behaviour belong together. Keeping tui-launcher's picker on the same key would make two nodes claim one key with different behaviour, which the model should not hold. The code today (internal/api/rex_install.go:37, launch_key = "ctrl+a a" bound to rhiza_launch) is the pre-change state this change set replaces, so it is not evidence against the answer. Not decided: whether the tui-launcher intent is retired or kept for another path; that is left for the owner.

## statement

On ctrl+a>a the native change set supersedes the tui-launcher picker: the key enters the rhiza key mode with one key per app, and no quick-launch layer lists the apps in a modal. tui-launcher's picker is not kept on that key; the key stays configurable in rhiza.lua, with ctrl+a>a as the default.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m644m1-r3-f2](../question/review-01m644m1-r3-f2.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
