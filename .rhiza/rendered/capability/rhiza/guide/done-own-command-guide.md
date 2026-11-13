# done-own-command-guide

capability · `01M76XY3S1VF65NT8YXQKZ4VZC` · parent [guide](../guide.md)

## acceptance

go test -trimpath ./internal/guide -run 'TestBuilderText|TestVerifierText|TestReviewText|TestSharedPhrasesInExactlyOneFile' asserts the shared session-done block holds the own-command rule, that the builder and verifier topics carry it and that review.md holds it; go test -trimpath ./internal/workflow asserts no prompt in internal/workflow/builtin/*.toml contains '; rhiza session done' or '&& rhiza session done' and every prompt that names rhiza session done also says it is the very last action.

## statement

The builder, verifier and review guide topics (the builder and verifier through the shared session-done block) and every builtin workflow prompt that ends with rhiza session done say to run it as its own command, only after the recording command (rhiza run findings, rhiza run verdict, the Evidence write or the commit) succeeded, never chained after it with ; or &&.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [done-own-command-guide-verified](../../../evidence/done-own-command-guide-verified.md)
- verifies [done-own-command-guide-qa-1](../../../evidence/done-own-command-guide-qa-1.md)
- verifies [done-own-command-guide-verified-2](../../../evidence/done-own-command-guide-verified-2.md)
- verifies [pr-495-done-needs-a-record](../../../evidence/pr-495-done-needs-a-record.md)
- verifies [pr-521-build-role-topics](../../../evidence/pr-521-build-role-topics.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
