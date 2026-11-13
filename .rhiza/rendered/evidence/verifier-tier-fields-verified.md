# verifier-tier-fields-verified

evidence · `01M79FCRJHQ14961MYN5WF0XS7` · actor [claude-code-57cdad9f](../actor/claude-code-57cdad9f.md)

## result

PASS

## statement

At d357c3a0f (detached checkout) I ran go test -trimpath ./internal/workflow/... (ok), go test -trimpath ./cmd/rhiza -run 'TestScripts/(workflow-describe|describe-golden)$' (ok) and go test -trimpath ./internal/api/... (ok). TestVerifierTierLoads parses a leaves and a checks verifier with models, low_risk_models and risky_paths, and confirms a role setting neither has nil fields. TestVerifierTierRefused refuses both keys on a builder, fixer, critic, check and agent step, and refuses a blank model, an empty glob and a malformed glob. TestRisky matches **/testdata/**, *_test.go and .github/** against sample paths, positive and negative. workflow-describe.txtar checks both fields in text and --json, and describe-golden carries the new shape. Reading the diffs (333601a65 + d357c3a0f) showed risk.go matching, parse and validation in workflow.go, and describe output in internal/api/workflow_describe.go. Both attempt-1 findings are fixed: a scratch build's `workflow describe` on a check setting risky_paths now says 'check: low_risk_models and risky_paths are only for the verifier role; a check cannot have them, set them under [step.review.verifier]', and on an agent step it says 'belong to a leaves or checks step's verifier role; an agent step has no verifier and cannot have them'. The tests assert both messages.

## Links

### Out

- verifies [verifier-tier-fields](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-fields.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
