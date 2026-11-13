# verifier-tier-choice-qa-1

evidence · `01M79GKZQBTVVJP5REM9R3W1SM` · actor [claude-code-052de74d](../actor/claude-code-052de74d.md)

## result

PASS

## statement

QA at dde99439c with a scratch build (go build -trimpath -o $TMPDIR/rhiza-qa ./cmd/rhiza), in a scratch Rhiza root and repo with fake rex and fake claude scripts (the fake claude logs role, leaf, attempt and --model; builders commit one file per leaf, verifiers write Evidence), run under env -i. rhiza run start --foreground --workflow build on three chained leaves: plain (commits internal/x/plain.go), risky (.github/workflows/ci.yml), testy (internal/x/x_test.go). Run succeeded; the journal's leaf-session models were plain verifier sonnet, risky verifier opus, testy verifier opus, builders sonnet. With plain's first verifier failing: plain builder 2 opus, plain verifier attempt 2 opus (retry stays strong), others as before, run succeeded. With a project .rhiza/workflows/build.toml overriding the verifier to models [opus] only: all three verifiers opus (no tiering, as before). With an override of low_risk_models [haiku], risky_paths [internal/x/risky.go, **/*.yml]: plain haiku, risky opus, testy haiku. Review's fixed-leaf verifier was not driven hands-on; go test -trimpath ./internal/run -run 'TestVerifierTier|TestVerifierRisky' ok.

## Links

### Out

- verifies [verifier-tier-choice](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-choice.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
