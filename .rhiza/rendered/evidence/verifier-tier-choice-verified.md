# verifier-tier-choice-verified

evidence · `01M79G8M6MQ79NZGW9VRD5QXE7` · actor [claude-code-0879b11d](../actor/claude-code-0879b11d.md)

## result

PASS

## statement

Verifier (attempt 2) at ea1121b28. Read the diffs of cbbf0ef76 and ea1121b28: startLeaf in internal/run/leaves.go switches a verifier to low_risk_models only when verifierRisky is false; verifierRisky has no agent call and returns risky when the round attempt is past 1, the round is past 1, the review cycle is past 1, the runner cannot list changed files, the builder's/fixer's journalled start HEAD or the commit is missing, ChangedFiles errors, or a changed file matches risky_paths. stepRunner.ChangedFiles runs git diff --name-only -z. The spec's Model is what the leaf-session entry journals. Ran go test -trimpath ./internal/run/ -run 'TestVerifierTier|TestVerifierRisky' -v: all pass. With a role setting low_risk_models = [sonnet], models = [opus], risky_paths = [.github/**], the tests read the journal's verifier leaf-session models: sonnet for a first attempt changing only internal/x/x.go (diff range start HEAD..leaf commit checked), opus when .github/workflows/ci.yml changes, opus when the diff errors, opus on a role without low_risk_models; attempt 2 journals leaf-a 1 sonnet then leaf-a 2 opus (TestVerifierTierRetry); a retry in a fresh round journals opus for the round-2 verifier and sonnet for the other leaves' low-risk first attempts (TestVerifierTierFreshRound); review cycle 2 of a checks step journals opus (TestVerifierTierReviewCycle). This resolves attempt 1's three findings. Also ran go test -trimpath ./internal/run/... ./internal/api/ ./internal/workflow/: ok.

## Links

### Out

- verifies [verifier-tier-choice](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-choice.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
