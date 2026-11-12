# done-refuses-unrecorded-check

capability · `01M76XY3QZ8A5JTZE1V4EMAZMN` · parent [agent-step-done-signal](../agent-step-done-signal.md)

## acceptance

Go tests in internal/api (go test -trimpath ./internal/api -run SessionDone) assert: a check session with no round for its check and cycle is refused as a conflict naming rhiza run findings; after rhiza run findings --check with {"findings":[]} it succeeds; a critic session with one blocker unjudged is refused naming rhiza run verdict, and succeeds once every blocker and major has a verdict; a fixer session is never refused for this; the existing plan reviewer case still passes.

## statement

rhiza session done from a review check session refuses until a findings round (empty for a clean check) is journaled for its check and cycle, and from a critic session until a verdict is journaled for every blocker and major finding of its check's latest round, as plan reviewers already are.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [done-refuses-unrecorded-check-verified](../../../../../evidence/done-refuses-unrecorded-check-verified.md)
- verifies [done-refuses-unrecorded-check-qa-1](../../../../../evidence/done-refuses-unrecorded-check-qa-1.md)
- verifies [pr-495-done-needs-a-record](../../../../../evidence/pr-495-done-needs-a-record.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
