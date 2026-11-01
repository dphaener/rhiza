# check-warns-overlaps-verified

evidence · `01M6A96JJKHNXNEXK580ES4NNX` · actor [claude-code-eea0e42d](../actor/claude-code-eea0e42d.md)

## result

PASS

## statement

At 8f186a10: go test ./internal/check (new TestRun 'overlap' case: two proposed change sets set one capability's field; one warn in-flight finding naming both sets and the node, ChangeSet A, Next 'rhiza slice --change A --view what-if') passes; go test ./cmd/rhiza -run TestScripts/check-overlaps (new script: two proposed sets set statement of one node; rhiza check prints ok and exactly one 'warn in-flight ... change sets X and Y both write field:statement of N' with next rhiza slice --change X --view what-if; --json ok true, no fail) passes; full go test ./... passes, so no existing script output needed updating; diff read: check.go adds model.Overlaps findings under RuleInFlight. make check's self-check fails only on stale .rhiza/rendered from prior commit c0fa1d22 (same failure at HEAD~1; leaf touches no .rhiza file).

## Links

### Out

- verifies [check-warns-overlaps](../capability/rhiza/integrity/check/check-warns-overlaps.md)

## History

- `01M5SWTNWB9GEJ5AA8R7K867KB` Genesis detects conflicting decisions across change sets in flight and surfaces them to the owner.
