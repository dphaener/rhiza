# api-timings-seed-after-fast-and-opt-in

decision · `01M7FVAQVZNGDMEP788DRYGQBJ`

## reasoning

The finding holds: api-test-timings has no depends-on link, so a conductor may seed the file before link-endpoint-kinds-fast cuts TestLinkEndpointKinds from about 81 s to under 15 s and before every-budget-opt-in adds the two budget tests to BUDGET_TESTS. api-shards-weighted would then place tests by stale weights and miss its CI measure (api jobs within 20 s) by about 66 s, and nothing refreshes the file afterwards. The plan already settles the same staleness on the script side by ordering (script-timings-measured depends on fake-rex-no-dead-wait), and api-shards-weighted-like-scripts says the api side follows the scripts, so ordering is the consistent fix. It is also the cheapest: two depends-on links, no new leaf or CI refresh step, and the acceptance's seed command (-skip BUDGET_TESTS) stays as written.

## statement

api-test-timings is built after link-endpoint-kinds-fast and every-budget-opt-in: it depends on both, so its seed of internal/api/testdata/test-timings.tsv is taken from a go test -json run of internal/api with TestLinkEndpointKinds already fast and with TestInboxGateQuestionsBudget and TestRootIdlePollBudget already named in BUDGET_TESTS (and so skipped and left out of the file).

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7frv7b35b91y35xa8vxzpc0-r3-f1](../question/review-01m7frv7b35b91y35xa8vxzpc0-r3-f1.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
