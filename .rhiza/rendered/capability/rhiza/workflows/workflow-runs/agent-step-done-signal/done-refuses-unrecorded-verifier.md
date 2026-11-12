# done-refuses-unrecorded-verifier

capability · `01M76XY3PYRWRTE99N44HPC2G1` · parent [agent-step-done-signal](../agent-step-done-signal.md)

## acceptance

Go tests in internal/api (go test -trimpath ./internal/api -run SessionDone) start a verifier session on a leaf attempt and assert: with nothing recorded, session done exits as a conflict naming rhiza run findings and writes no done mark; after rhiza run findings on that attempt it succeeds; with Evidence by the verifier's actor it succeeds; Evidence by another actor alone is still refused; a builder session is never refused for this.

## statement

rhiza session done from a build leaf's verifier session refuses, naming rhiza run findings and the Evidence recipe, while that leaf's attempt has neither journaled findings nor Evidence written by the verifier's actor, so a verifier cannot end an attempt with nothing recorded.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [done-refuses-unrecorded-verifier-verified](../../../../../evidence/done-refuses-unrecorded-verifier-verified.md)
- verifies [done-refuses-unrecorded-verifier-qa-1](../../../../../evidence/done-refuses-unrecorded-verifier-qa-1.md)
- verifies [pr-495-done-needs-a-record](../../../../../evidence/pr-495-done-needs-a-record.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
