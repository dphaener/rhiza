# tui-many-mini-apps

decision · `01M62JMM4GYAAGBWRJGMP7X19P`

## quote

It should launch from a rex quick actions layer. This will require us to have an install command that installs a lua script in the rex config. We should be able to do something like ctrl-a>r to open a rex specific quick launch rex layer and then select from there which rex "mini" app we want to open. This isn't just a single one size fits all tui, it's many mini apps that sort of wrap around Rex and give you small bits of functionality rather than one large TUI that is always open.

## statement

The TUI is many small apps launched from a Rex quick-launch layer, installed by a command that writes a Lua script into the Rex config, not one large TUI that is always open

## Links

### Out

- derived-from [tui-planning-conversation](../source/tui-planning-conversation.md)

### In

- decided-by [rhiza-tui](../intent/rhiza-tui.md)

## History

- `01M62J0ED188JTTR7R9709Y989` Rhiza has a terminal UI built on Rex, with rexuikit, in which common rhiza commands run as tiny apps in modals or pop-ups, alongside a view of the model.
