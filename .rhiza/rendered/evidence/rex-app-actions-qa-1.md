# rex-app-actions-qa-1

evidence · `01M648A0A6SZAZWBEE26MWVQMZ` · actor [claude-code-fe283bac](../actor/claude-code-fe283bac.md)

## result

FAIL

## statement

Built rhiza at 8182d27, ran rhiza rex install into a scratch REX_CONFIG (fake rex for the reload), then loaded it with the real rex config check: rhiza.lua has 16 rhiza_app_* actions, one per app, titled by name and running rhiza tui app <name>, but real Rex rejects 7 of them (change-show, record-add, record-link, record-retire, record-set, run-status, session-ls): 'action name "rhiza_app_session-ls" must match [a-z0-9_]+', and config check exits 1 with 'the configuration has errors'. Only 9 apps (plus rhiza_launch) become actions.

## Links

### Out

- verifies [rex-app-actions](../capability/rhiza/rex-app-actions.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
