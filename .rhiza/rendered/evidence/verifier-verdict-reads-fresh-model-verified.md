# verifier-verdict-reads-fresh-model-verified

evidence · `01M796HG09F7H0CN6BM7XZZSQY` · actor [claude-code-cb564ef2](../actor/claude-code-cb564ef2.md)

## result

PASS

## statement

At 8e90ac317, in throwaway worktrees (not committed): with pollDebounce raised to 10s as the injected delay and r.settle() removed from stepRunner.Evidence, go test -trimpath ./cmd/rhiza -run 'TestScripts/headless-run-agent-tabs$' failed twice with 'verifier ended with neither Evidence nor findings' (the flake); with the same delay and the fix it passed 3 consecutive runs. Without the fix, go test ./internal/api -run TestEvidenceReadsPastSessionWrites failed (Evidence after the session's writes = false; want true). Without the delay, at the leaf commit: headless-run-agent-tabs -count=5 ok, go test -trimpath ./internal/api/... ok, go test -trimpath ./internal/run/... ok. The script has no skip or quarantine. Diff read: settle() under modelMu sets catchUpNext, refreshDue honours it, catchUpModel clears it after a catch-up.

## Links

### Out

- verifies [verifier-verdict-reads-fresh-model](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaves-step/verifier-verdict-reads-fresh-model.md)

## History

- `01M7942GRFNQ82YV0V8F4A9VNA` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172; reproduce it locally, repeatedly, then fix it or quarantine it
