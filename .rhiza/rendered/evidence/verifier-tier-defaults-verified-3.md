# verifier-tier-defaults-verified-3

evidence · `01M79JBQP6REMSW3KTEGH8SKQ2` · actor [claude-code-e467b1d7](../actor/claude-code-e467b1d7.md)

## result

PASS

## statement

At 28fe31869: go test -trimpath ./internal/workflow/... passes, and TestBuiltinVerifierTierDefaults parses the built-in build and review workflows and finds models [opus], low_risk_models [sonnet] and the full risky_paths list (test infra, hooks, .github, workflow files, migrations, *retry*) on both verifier roles. go test -trimpath ./cmd/rhiza -run 'TestScripts/(workflow-describe|build-workflow)$' passes; workflow-describe asserts 'role verifier low_risk_models [sonnet]' and the risky_paths line for `rhiza workflow describe build`. make self-check passes (rc=0) in a clean detached worktree at 28fe31869 under env -i. In the live run worktree it fails only on this verifier's own untracked actor render. Read build.toml and review.toml: no new rhiza.* setting. The earlier finding is fixed: build-workflow-file's statement now reads 'verifier opus for risky leaves, sonnet for low-risk ones'.

## Links

### Out

- verifies [verifier-tier-defaults](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-defaults.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
