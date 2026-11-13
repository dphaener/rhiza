# verifier-tier-defaults-verified

evidence · `01M79FTBX76W9J224BVQM7XAMM` · actor [claude-code-ba054218](../actor/claude-code-ba054218.md)

## result

PASS

## statement

At 6af6a212a: go test -trimpath -count=1 ./internal/workflow/... ok (TestBuiltinVerifierTierDefaults parses the built-in build and review workflows and asserts models [opus], low_risk_models [sonnet] and the twelve risky_paths globs on both verifier roles); go test -trimpath ./cmd/rhiza -run 'TestScripts/(workflow-describe|describe-golden)$' ok; a scratch build's `rhiza workflow describe build` and `describe review` both print role verifier models [opus], low_risk_models [sonnet] and risky_paths [*_test.go, **/*_test.go, **/testdata/**, **/main_test.go, .githooks/**, **/hooks/**, .github/**, .rhiza/workflows/**, internal/workflow/builtin/**, **/migrations/**, **/migrate/**, **/*retry*]. The diff only touches build.toml, review.toml, the test and the testscript; it adds no rhiza.* setting. make self-check in a clean detached checkout with env -i fails only on rendered-file staleness (actor/claude-code-57cdad9f.md missing, index.md differs). The identical failure is present at parent dee8a4784 (the run's evidence commit), so it is run-machinery render drift that `rhiza render` at ship will fix. Every other self-check step passes.

## Links

### Out

- verifies [verifier-tier-defaults](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-defaults.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
