# done-own-command-guide

capability · `01M76XY3S1VF65NT8YXQKZ4VZC` · parent [guide](../guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes: TestBuilderText, TestVerifierText and TestSharedPhrasesInExactlyOneFile assert the shared session-done block holds the own-command rule and the builder and verifier topics carry it, each of the nine review role topic tests (TestCheckComplianceText, TestCheckQualityText, TestCheckQAText, TestCheckModelText, TestCriticText, TestFixerText, TestCIFixText, TestGateFixText, TestConflictFixText) asserts its expanded topic carries the session-done block, and TestReviewRoleText finds the rule (run rhiza session done as its own command, after the recording command succeeded, never with ; or &&) in the role topics; go test -trimpath ./internal/workflow asserts no prompt in internal/workflow/builtin/*.toml contains '; rhiza session done' or '&& rhiza session done' and every prompt that names rhiza session done says it is the very last action (one-line /rhiza-<role> prompts that no longer name it pass).

## statement

The builder and verifier guide topics and every review role topic (check-compliance, check-quality, check-qa, check-model, critic, fixer, ci-fix, gate-fix and conflict-fix), each through the shared session-done block, and every builtin workflow prompt that still names rhiza session done say to run it as its own command, only after the recording command (rhiza run findings, rhiza run verdict, the Evidence write or the commit) succeeded, never chained after it with ; or &&; the review index does not hold the rule, and one-line /rhiza-<role> prompts leave it to their topic.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../decision/efficient-cycle-milestone-shape.md)

### In

- blocks [review-01m79a41352xbz4fc2fjt5903j-r3-f5](../../../question/review-01m79a41352xbz4fc2fjt5903j-r3-f5.md)
- verifies [done-own-command-guide-verified](../../../evidence/done-own-command-guide-verified.md)
- verifies [done-own-command-guide-qa-1](../../../evidence/done-own-command-guide-qa-1.md)
- verifies [done-own-command-guide-verified-2](../../../evidence/done-own-command-guide-verified-2.md)
- verifies [pr-495-done-needs-a-record](../../../evidence/pr-495-done-needs-a-record.md)
- verifies [pr-521-build-role-topics](../../../evidence/pr-521-build-role-topics.md)
- verifies [pr-537-review-role-topics](../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
