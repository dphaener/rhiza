# idle-poll-budget-times-idle-poll-qa-1

evidence · `01M7BDY5D18RQ5FAX3CG3D6YSM` · actor [claude-code-9d795993](../actor/claude-code-9d795993.md)

## result

PASS

## statement

QA cycle 1 at a445e38bc: built rhiza fresh to a scratch path (go build ./cmd/rhiza, runs); go vet ./internal/api ok. Ran go test -trimpath -count=1 -p 1 -run '^TestReadBudget$' ./internal/api with RHIZA_BUDGET_LOG 3 times one after another: all ok (287s/360s/175s); driver poll idle medians 2.076/2.563/2.413 ms vs 5 ms budget (cpu pressure 44/87/48%), after-an-unrelated-entry medians 5.298/7.193/5.59 ms vs 50 ms. go test -trimpath -count=1 -run '^(TestDriverPollIsCheap|TestEvidenceReadsPastSessionWrites)$' ./internal/api ok. Diff touches only internal/api/run_model_cache_test.go: the idle measure passes leaf "" so driverFixture.poll skips Evidence (line 47 guard), warm-up and unrelated measures keep driver-leaf-a, budgets stay 5 ms and 50 ms, a comment explains why; stepRunner.Evidence and settle() untouched.

## Links

### Out

- verifies [idle-poll-budget-times-idle-poll](../capability/rhiza/orchestration/daemon/idle-poll-budget-times-idle-poll.md)

## History

- `01M7B7P8D389GB7N94YPTXAQCS` Fix flaky test TestReadBudget: failed then passed on a rerun in PR #567 at head c5cc453c8b33b60fd6078c0aeaf5e1331ce3239f; reproduce it locally, repeatedly, then fix it or quarantine it
