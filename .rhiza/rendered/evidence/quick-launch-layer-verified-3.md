# quick-launch-layer-verified-3

evidence · `01M6316S2HEPZXFA34MV775MN9` · actor [claude-code-b75d8bbf](../actor/claude-code-b75d8bbf.md)

## result

PASS

## statement

Verified commit 3f3f1f8 (attempt 3, the fix for QA cycle 2's major finding that the layer started in $HOME). Read the diff: the rhiza_launch action now passes cwd = block_cwd(ctx.block_id), which reads the invoking terminal's foreground (else child) cwd through rex.block.call(TERM, "process"), a documented terminal method. go test ./... passes, including the fake-Rex tests for binding the launch key, listing every app by name, and opening the chosen app's modal; fmt-check and vet pass. By hand, with real Rex: built the binary to a scratch path, ran rhiza rex install into a scratch REX_CONFIG with rex stubbed for the reload, and made a throwaway Rex session with -c set to this checkout. The process method gave foreground.cwd as the checkout. From /tmp I ran rex do rhiza.lua --action rhiza_launch -s/-b on that session's block. The layer opened listing ready, roadmap, show, slice and why; pressing enter on ready opened a 'ready' modal showing the repo's change_sets JSON, not 'no Rhiza repository'. Note: make check's self-check fails only on stale .rhiza/rendered actor/index files from the run's model writes; this commit changes no model or rendered file.

## Links

### Out

- verifies [quick-launch-layer](../capability/rhiza/quick-launch-layer.md)

## History

- `01M62JMSW4FH8RAWSQXHC7ZPND` Rex quick-launch layer and install command, with roadmap, ready, show, why and slice as apps
