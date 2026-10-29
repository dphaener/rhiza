# quick-launch-layer-qa-3

evidence · `01M631J5TPJ1T3EC3X7BCSTN4Z` · actor [claude-code-2ac42aad](../actor/claude-code-2ac42aad.md)

## result

PASS

## statement

QA at d02d9ce (cycle 3) in real Rex. Built rhiza to a scratch dir, ran rhiza rex install into a scratch REX_CONFIG, then ran a throwaway rex server run on a scratch socket with --config at that init.lua (REX_SERVER pointed at itself, RHIZA_* unset). rex keymap on it lists 'ctrl+a a rhiza_launch rhiza.lua:44' beside the owner's 'ctrl+a r client.config.reload', and rex actions lists rhiza_launch. Made a session with -c at the scratch dir, then cd'd its shell into this checkout, and ran rex do rhiza_launch -s/-b on that shell block: a 'Rhiza apps' layer opened whose process cwd was the checkout (not $HOME nor the session's start dir), listing ready, roadmap, show, slice, why with summaries. send-key enter on ready closed it and opened a 'ready' block with this repo's change_sets JSON and the 'ready · j/k scroll · q quit' footer; the cycle-2 'no Rhiza repository' failure is fixed. Gap: the ctrl+a a keypress itself was not pressed in a real Rex.app client; when the action is called with no block in its context (rex.invoke from a script) the layer starts in the session's start dir.

## Links

### Out

- verifies [quick-launch-layer](../capability/rhiza/quick-launch-layer.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
