# idle-poll-budget-times-idle-poll-verified

evidence · `01M7BCSBHFBY80D19HX5BGAJFS` · actor [claude-code-5106d849](../actor/claude-code-5106d849.md)

## result

PASS

## statement

Verified commit c765796d3. go test -trimpath -count=1 -p 1 -run '^TestReadBudget$' ./internal/api with RHIZA_BUDGET_LOG, 3 runs one after another: at the parent 35d4cdb4d exit 0 each, driver poll idle medians 7.335/5.64/3.941 ms (cpu pressure 69/68/4%, the first two inconclusive over 5 ms); at the leaf exit 0 each, idle medians 2.164/2.442/2.399 ms (pressure 38/57/64%), so the middle value fell from 5.64 to 2.399 ms. The 'after an unrelated entry' medians were 4.4-6.7 ms against 50 ms. go test -trimpath -count=1 -run '^(TestDriverPollIsCheap|TestEvidenceReadsPastSessionWrites)$' ./internal/api: exit 0. rhiza ship generate: exit 0. rhiza render: exit 0. In the diff, the idle measure polls with leaf "" so f.poll skips Evidence, and a comment says why. The warm-up polls and the unrelated measure still pass driver-leaf-a. idlePollCPUBudget is still 5 ms and unrelatedPollCPUBudget 50 ms. The only code file changed is internal/api/run_model_cache_test.go, so stepRunner.Evidence and settle() are unchanged; the other files are the builder's staged actor record and its render. No TODOs or stubs.

## Links

### Out

- verifies [idle-poll-budget-times-idle-poll](../capability/rhiza/orchestration/daemon/idle-poll-budget-times-idle-poll.md)

## History

- `01M7B7P8D389GB7N94YPTXAQCS` Fix flaky test TestReadBudget: failed then passed on a rerun in PR #567 at head c5cc453c8b33b60fd6078c0aeaf5e1331ce3239f; reproduce it locally, repeatedly, then fix it or quarantine it
