# tui-launch-key

decision · `01M62VBX0YHEBNPTRZK2C8CJ7R`

## reasoning

The owner asked for 'something like ctrl-a>r'. ~/.config/rex/init.lua binds ctrl+a r to client.config.reload, and taking it would break that binding. ctrl+a a ('apps') is unbound, sits under the same prefix, and the setting lets the owner pick another key.

## statement

The quick-launch layer opens on ctrl+a a by default, not ctrl+a r, because the owner's init.lua binds ctrl+a r to client.config.reload. The key is a setting in rhiza.lua that a reinstall keeps.

## wording

the orchestrator's recommended option, taken under the owner's standing word: "i feel that I've made enough high level decisions" and "I'm just accepting the recommended answers anyway"

## Links

### In

- decided-by quick-launch-layer (retired)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
