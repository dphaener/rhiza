# review-01m644m1-r3-f2

question · `01M6457MCVYTZCTTDE6E1HSPM4`

## statement

Plan review round 3 left finding f2 (contradiction on tui-launcher): The tui-launcher intent (01M62JMJ6AYPDDEEFNPG9WG3E0) says ctrl-a>a opens a Rex quick-launch layer that lists every rhiza app and opens the chosen one in a modal. The plan's native-rex-prefix-configurable decision and rex-launch-key-setting capability make the same ctrl+a>a enter a rhiza key mode with one key per app, with no picker in between. The two nodes disagree on what ctrl+a>a does. native-intent-served (01M644QZSWE52YYK67FX8B54SK) asks whether this supersedes tui-launcher, but no node records the answer yet, so the conflict stays in the model. The critic confirmed it: tui-launcher intent (01M62JMJ6AYPDDEEFNPG9WG3E0) says ctrl-a>a opens a quick-launch layer listing every rhiza app, with a picker and a modal. The plan's native-rex-prefix-configurable decision and rex-key-mode capability make the same ctrl+a>a enter a key mode with one key per app and no picker. The code agrees with the intent, not the plan: internal/api/rex_install.go:37 writes launch_key = "ctrl+a a" bound to rhiza_launch, the all-apps picker. The question native-intent-served (01M644QZSWE52YYK67FX8B54SK) names the conflict but no node records how it is resolved, so the contradiction stands in the model. How should it be settled?

## Links

### In

- answers [native-supersedes-tui-launcher](../decision/native-supersedes-tui-launcher.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
