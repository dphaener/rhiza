# plan-autonomous-journal-order-ignores-tokens-qa-1

evidence · `01M7RGY2BKAHCWY9E8PHAF7XFY` · actor [claude-code-8045855a](../actor/claude-code-8045855a.md)

## result

PASS

## statement

QA at 60d5352f0, each go test in a systemd-run scope with -p 2: plan-workflow-autonomous -count=3 ok; with the two scratch delays (engine.go Sleep 1500ms after the step log, sleep 0.8 after the fake planner's session done) the fixed script passed 3/3 while main f7313222e's failed 3/3 at txtar:68 with 'unblocked tokens tokens ended tokens entered'; tokens step changed to nosuch failed at txtar:72 and no tokens entries at all failed at txtar:71; the -left and -unapplied scripts pass; all scratch edits discarded

## Links

### Out

- verifies [plan-autonomous-journal-order-ignores-tokens](../capability/plan-autonomous-journal-order-ignores-tokens.md)

## History

- `01M7RFDAZ53YZ3775D7SV70Z4G` Fix flaky test TestScripts/plan-workflow-autonomous: failed then passed on a rerun in PR #655 at head 1be14770b82366ab1535a1b38a36a890e8b75939; reproduce it locally, repeatedly, then fix it or quarantine it
