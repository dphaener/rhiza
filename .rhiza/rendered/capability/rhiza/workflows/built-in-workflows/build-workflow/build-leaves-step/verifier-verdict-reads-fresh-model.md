# verifier-verdict-reads-fresh-model

capability · `01M794A1BTCCFC0G1V1MPQ5V89` · parent [build-leaves-step](../build-leaves-step.md)

## acceptance

With a temporary delay injected at the race (not committed), go test -trimpath ./cmd/rhiza -run 'TestScripts/headless-run-agent-tabs$' fails before the fix and passes after it on 3 consecutive runs; a new regression test in internal/api (or internal/run) that writes a session's Evidence after the cache has kept an answer, inside the debounce window, fails without the fix and passes with it; without the delay go test -trimpath ./cmd/rhiza -run 'TestScripts/headless-run-agent-tabs$' -count=5 passes, as do go test -trimpath ./internal/api/... and go test -trimpath ./internal/run/... when those packages changed; the script is not quarantined, unless the race could not be reproduced or fixed, in which case it is skipped with a bug request filed via rhiza request --kind bug and the reason written as a note on this change set

## statement

A leaves step judges a session that has just ended (a verifier's Evidence, a builder's or fixer's leaf commit) against a model caught up past that session's writes, never an answer the run driver's model cache (internal/api/run_model_cache.go: readModel, catchUpModel, refreshDue, revalidate) kept from before them, so a verifier that wrote Evidence while another leaf's session wrote concurrently is judged passed, and headless-run-agent-tabs no longer flakes

## Links

### Out

- decided-by [verifier-evidence-read-fresh](../../../../../../decision/verifier-evidence-read-fresh.md)
- derived-from [headless-tabs-flake-plan](../../../../../../source/headless-tabs-flake-plan.md)

### In

- verifies [verifier-verdict-reads-fresh-model-verified](../../../../../../evidence/verifier-verdict-reads-fresh-model-verified.md)
- verifies [verifier-verdict-reads-fresh-model-qa-1](../../../../../../evidence/verifier-verdict-reads-fresh-model-qa-1.md)
- verifies [pr-516-fix-flaky-test](../../../../../../evidence/pr-516-fix-flaky-test.md)

## History

- `01M7942GRFNQ82YV0V8F4A9VNA` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172; reproduce it locally, repeatedly, then fix it or quarantine it
