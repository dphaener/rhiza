# verifier-tier-choice-verified-2

evidence · `01M79HR5BKXVQ1645NCB07VYPA` · actor [claude-code-5905f72b](../actor/claude-code-5905f72b.md)

## result

PASS

## statement

At 72436a158 in a detached checkout: go test -trimpath ./internal/run/ -run 'TestVerifierTier|TestVerifierRisky' -v passes all six tests, and go test -trimpath ./internal/run/... passes. The tests start verifiers from a role with models [opus], low_risk_models [sonnet], risky_paths [.github/**]. They journal sonnet for a low-risk first attempt touching only internal/x/x.go (TestVerifierTierChoice) and for review cycle 2's first fix on that file (TestVerifierTierReviewCycle). They journal opus for .github/workflows/ci.yml, an unreadable diff and a role without low_risk_models (TestVerifierTierChoice), for attempt 2 (TestVerifierTierRetry), for a fresh round (TestVerifierTierFreshRound) and for cycle 2's second fix try (TestVerifierTierReviewRetry). Cycle 3 is risky whatever the diff (TestVerifierRiskyRetryStates). The diff to verifierRisky and startLeaf matches: a checks step is a retry when the cycle is past 2, the try is not leafFixAttempt(cycle,1), or the round is past 1, and a low-risk launch takes the cheap ladder's first model. make self-check exits 0. build-workflow.txtar fails at :46 (it expects opus verifiers); that comes from the build.toml defaults and is already a confirmed review finding, not part of this leaf's acceptance.

## Links

### Out

- verifies [verifier-tier-choice](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-choice.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
