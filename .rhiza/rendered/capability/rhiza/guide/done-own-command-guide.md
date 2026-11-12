# done-own-command-guide

capability · `01M76XY3S1VF65NT8YXQKZ4VZC` · parent [guide](../guide.md)

## acceptance

go test -trimpath ./internal/guide -run 'Build|Review' asserts build.md and review.md each hold the own-command rule; go test -trimpath ./internal/workflow asserts no prompt in internal/workflow/builtin/*.toml contains '; rhiza session done' or '&& rhiza session done' and every prompt that names rhiza session done also says it is the very last action.

## statement

The build and review guide topics and every builtin workflow prompt that ends with rhiza session done say to run it as its own command, only after the recording command (rhiza run findings, rhiza run verdict, the Evidence write or the commit) succeeded, never chained after it with ; or &&.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [done-own-command-guide-verified](../../../evidence/done-own-command-guide-verified.md)
- verifies [done-own-command-guide-qa-1](../../../evidence/done-own-command-guide-qa-1.md)
- verifies [done-own-command-guide-verified-2](../../../evidence/done-own-command-guide-verified-2.md)
- verifies [pr-495-done-needs-a-record](../../../evidence/pr-495-done-needs-a-record.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
