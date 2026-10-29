# quick-launch-layer-qa-2

evidence · `01M630C3ZCV8W1AA33THRJ328X` · actor [claude-code-f3ca46cd](../actor/claude-code-f3ca46cd.md)

## result

FAIL

## statement

QA at c7a97c9 (cycle 2) in real Rex. Built rhiza to a scratch dir; rex new a throwaway session with --cwd at the repo running `rhiza tui launch`: the layer listed ready, roadmap, show, slice, why with summaries; down+enter opened a floating 'roadmap' block showing the roadmap JSON with the 'roadmap · j/k scroll · q quit' footer, so the cycle-1 'stream is not attached (409)' failure is fixed. Then ran the installed rhiza.lua's rhiza_launch action in real Rex with `rex do -s <session> rhiza.lua --action rhiza_launch` (rex.bind line stripped so the owner's keymap was untouched): it opened a focused floating 'Rhiza apps' layer at bounds 0.3,0.25,0.4,0.5, and enter on ready closed it and opened a 'ready' modal. But that modal showed 'rhiza: no Rhiza repository in /Users/darinhaener': a pwd probe through the same rex.session.new_layer call shows a block with no cwd starts in $HOME, and rhiza.lua passes no cwd, so the launcher (and the apps, which take its os.Getwd) run outside the repo. The ctrl+a a keypress itself was not pressed in a real client.

## Links

### Out

- verifies [quick-launch-layer](../capability/rhiza/quick-launch-layer.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
