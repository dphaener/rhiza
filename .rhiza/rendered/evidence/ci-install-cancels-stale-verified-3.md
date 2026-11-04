# ci-install-cancels-stale-verified-3

evidence · `01M6M3K0Y0AB565BNZHRGNSBG8` · actor [claude-code-40c277e0](../actor/claude-code-40c277e0.md)

## result

PASS

## statement

Verifier, attempt 5 at a8b12607: read internal/api/ci.go and cmd/rhiza/testdata/script/ci-install.txtar. The testscript's expected rhiza.yml (cmp against the written .github/workflows/rhiza.yml) holds concurrency group rhiza-${{ github.event.pull_request.number || github.sha }} (PR number, falling back to the commit for pushes to main) with cancel-in-progress: ${{ github.event_name == 'pull_request' }}, keeps on.push.branches [main] and the pull_request trigger, and keeps the base-branch, before-commit and first-push check steps. The leaf's acceptance and statement now name the github.sha fallback, resolving the review's model finding. go test ./cmd/rhiza -run TestScripts/ci-install passed; make fmt-check vet test passed (rc 0). make self-check fails only on stale .rhiza/rendered actor files and index.md from in-run actor records, unrelated to this leaf.

## Links

### Out

- verifies [ci-install-cancels-stale](../capability/rhiza/integrity/ci-install-cancels-stale.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
