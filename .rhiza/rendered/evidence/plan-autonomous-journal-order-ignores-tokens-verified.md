# plan-autonomous-journal-order-ignores-tokens-verified

evidence · `01M7RGM5PDQRW95B6EKVTTA66J` · actor [claude-code-b1daa646](../actor/claude-code-b1daa646.md)

## result

PASS

## statement

Verified at 60d5352f0 in a detached checkout, each go test inside systemd-run --user --scope -p MemoryMax=6G -p TasksMax=300 with -p 2 and RHIZA_/SHIP_/REX_ unset, one at a time. (1) go test -trimpath -p 2 -count=3 ./cmd/rhiza -run 'TestScripts/plan-workflow-autonomous$': ok. (2) Scratch delays (time.Sleep(1500ms) after r.log step in internal/run/engine.go:772, sleep 0.8 after the fake planner's rhiza session done): with main f7313222e's txtar, -count=3 failed 3 of 3 at txtar:68 with 'unblocked tokens tokens ended tokens entered'; with the fix, -count=3 ok. (3) Tokens assertion's step changed to nosuch in scratch: exit 1, no match at txtar:72; reverted. (4) go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/plan-workflow-autonomous' (with -left, -unapplied): ok. rhiza ship generate and rhiza render: exit 0, no diff. Diff read: only plan-workflow-autonomous.txtar changes (plus the builder's actor record); the order regex is unchanged except the tokens group, which now is filtered out with grep -vx tokens; no Go code, skip list or other script touched; all scratch edits reverted and the checkout removed.

## Links

### Out

- verifies [plan-autonomous-journal-order-ignores-tokens](../capability/plan-autonomous-journal-order-ignores-tokens.md)

## History

- `01M7RFDAZ53YZ3775D7SV70Z4G` Fix flaky test TestScripts/plan-workflow-autonomous: failed then passed on a rerun in PR #655 at head 1be14770b82366ab1535a1b38a36a890e8b75939; reproduce it locally, repeatedly, then fix it or quarantine it
