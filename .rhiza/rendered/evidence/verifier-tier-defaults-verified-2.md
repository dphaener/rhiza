# verifier-tier-defaults-verified-2

evidence · `01M79HW35A3G50VE4R0YV58F2N` · actor [claude-code-ac4909dd](../actor/claude-code-ac4909dd.md)

## result

PASS

## statement

Review cycle 2 re-verification at 5c9808a08. I read build.toml and review.toml: both verifiers keep models=["opus"] and add low_risk_models=["sonnet"] and the twelve risky_paths globs the leaf lists, and no new rhiza.* setting was added. The fix commit only updates build-workflow.txtar to expect sonnet verifiers for its low-risk leaves. go test -trimpath -count=1 ./internal/workflow passes, including TestBuiltinVerifierTierDefaults, which loads the build and review built-ins and checks the ladders and globs. go test -trimpath ./cmd/rhiza -run 'TestScripts/(workflow-describe|build-workflow|describe-golden)$' passes; workflow-describe checks the low_risk_models and risky_paths lines of `rhiza workflow describe build`. make self-check exits 0 in a clean detached worktree at 5c9808a08 with the run's RHIZA_* environment unset. In the run worktree it fails only because the verifier's own actor (RHIZA_ACTOR) has not been rendered.

## Links

### Out

- verifies [verifier-tier-defaults](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-defaults.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
