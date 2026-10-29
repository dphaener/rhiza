# rex-key-mode-qa-2

evidence · `01M64A7TYMN7XCGMTFVJE5RZHT` · actor [claude-code-afc79974](../actor/claude-code-afc79974.md)

## result

PASS

## statement

At 7b2fd58, loaded the installed rhiza.lua in real Rex's Lua (rex do with a probe script recording rex.bind and rex.mode): mode rhiza exclusive=true; rhiza/escape -> client.mode.exit; one key per app, all distinct (a accept, b roadmap, c change-show, d record-retire, e record-add, f slice, h show, i session-ls, l land, m model, o record-link, p propose, r ready, s record-set, u run-status, w why), each bound to that app's rhiza_app_* action; ctrl+a>a -> client.mode.enter {name=rhiza, once=true}. With settings keys = { roadmap = "x", why = "" }, roadmap moved to rhiza/x and why kept its default w; a reinstall kept the override. Key presses themselves were not sent, since a headless session has no client.

## Links

### Out

- verifies [rex-key-mode](../capability/rhiza/rex-key-mode.md)

## History

- `01M643X7CEAQP16925KTAM06GT` Rhiza's apps open natively in Rex, as forge's quick actions do: rhiza rex install writes a Lua script that registers each app as its own Rex action (in the command palette) and a rhiza key mode with one key per app, and each action opens its app straight away in a layer sized to it, with no picker in between.
