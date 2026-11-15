# review-01m7frv7b35b91y35xa8vxzpc0-r3-f1

question · `01M7FV95ZCQ3ZJ0EJD59XEQF33`

## statement

Plan review round 3 left finding f1 (contradiction on api-test-timings): api-test-timings seeds internal/api/testdata/test-timings.tsv from one local go test -json run of internal/api, but it has no depends-on on link-endpoint-kinds-fast or every-budget-opt-in, so the conductor may build it first. Then the seed records TestLinkEndpointKinds at its old ~81 s, and TestInboxGateQuestionsBudget and TestRootIdlePollBudget as ordinary tests, because BUDGET_TESTS does not name them yet. api-shards-weighted would then place tests by those stale weights: about 81 s for a test link-endpoint-kinds-fast cuts to under 15 s, which defeats its own measure (the two api jobs within 20 s of each other). The plan handles the same staleness for the scripts in script-timings-from-ci-then-local, where script-timings-measured depends on fake-rex-no-dead-wait and re-times after it. The api side lacks that ordering: api-test-timings should depend on link-endpoint-kinds-fast and every-budget-opt-in, or its seed should be taken after both land. The critic confirmed it: Checked with rhiza show: api-test-timings has no depends-on link at all (only decided-by api-shards-weighted-like-scripts, and api-shards-weighted depends on it). link-endpoint-kinds-fast and every-budget-opt-in have no ordering links either, while script-timings-measured does depend on fake-rex-no-dead-wait. So the conductor may build api-test-timings first. Its acceptance seeds the checked-in file from one local go test -json run of internal/api, which would then record TestLinkEndpointKinds at about 81 s (research: 81 s alone, cut to under 15 s by link-endpoint-kinds-fast) and would time TestInboxGateQuestionsBudget and TestRootIdlePollBudget as ordinary tests, because BUDGET_TESTS does not name them yet. api-shards-weighted would place tests by those stale weights. Its file-only check (totals within 15% of the mean) still passes, but its CI measure (the two api jobs within 20 s of each other) would miss by about 66 s. Nothing in the plan refreshes the file from CI afterwards. This is a real defect in the plan's ordering that defeats a leaf's measure, so major stands. How should it be settled?

## Links

### Out

- blocks [api-test-timings](../capability/api-test-timings.md)

### In

- answers [api-timings-seed-after-fast-and-opt-in](../decision/api-timings-seed-after-fast-and-opt-in.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
