# verifier-tier-defaults-qa-1

evidence · `01M79GKZS1Z75X718ZYQE9DGBN` · actor [claude-code-052de74d](../actor/claude-code-052de74d.md)

## result

PASS

## statement

QA at dde99439c with a scratch build (go build -trimpath -o $TMPDIR/rhiza-qa ./cmd/rhiza), in a scratch Rhiza root and repo with fake rex and fake claude scripts (the fake claude logs role, leaf, attempt and --model; builders commit one file per leaf, verifiers write Evidence), run under env -i. rhiza workflow describe build and describe review each print role verifier models [opus], low_risk_models [sonnet] and the twelve risky_paths globs (*_test.go, **/*_test.go, **/testdata/**, **/main_test.go, .githooks/**, **/hooks/**, .github/**, .rhiza/workflows/**, internal/workflow/builtin/**, **/migrations/**, **/migrate/**, **/*retry*). A run of the built-in build workflow verified a leaf changing internal/x/plain.go on sonnet and leaves changing .github/workflows/ci.yml and internal/x/x_test.go on opus. go test -trimpath ./internal/workflow/... ok. Note: cmd/rhiza testscript build-workflow still expects opus verifiers for its low-risk leaves and now fails (passes at the merge-base), recorded as a finding.

## Links

### Out

- verifies [verifier-tier-defaults](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-defaults.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
