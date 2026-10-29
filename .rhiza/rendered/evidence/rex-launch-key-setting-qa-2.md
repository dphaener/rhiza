# rex-launch-key-setting-qa-2

evidence · `01M64A7VKJY4AQ58E4RB07PETJ` · actor [claude-code-afc79974](../actor/claude-code-afc79974.md)

## result

PASS

## statement

At 7b2fd58, a first install into a scratch REX_CONFIG wrote launch_key = "ctrl+a>a" and keys = {} in the settings block. Edited it to launch_key = "ctrl+b  r" and a keys map, reinstalled: --json reported written [] and unchanged [rhiza.lua, init.lua], the settings block was byte-identical (cmp), and the probe showed ctrl+b>r -> client.mode.enter. Reinstalling over a copy of the owner's current ~/.config/rex/rhiza.lua (launch_key = "ctrl+a a", bound to rhiza_launch) rewrote only rhiza.lua, kept that settings block, removed rhiza_launch, and bound ctrl+a>a to enter the mode; rex config check on it reports no diagnostics.

## Links

### Out

- verifies [rex-launch-key-setting](../capability/rhiza/rex-launch-key-setting.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
