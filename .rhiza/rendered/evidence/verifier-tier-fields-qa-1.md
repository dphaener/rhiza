# verifier-tier-fields-qa-1

evidence · `01M79GKZTFPNZJBH1D0KFE3QSG` · actor [claude-code-052de74d](../actor/claude-code-052de74d.md)

## result

PASS

## statement

QA at dde99439c with a scratch build (go build -trimpath -o $TMPDIR/rhiza-qa ./cmd/rhiza), in a scratch Rhiza root and repo with fake rex and fake claude scripts (the fake claude logs role, leaf, attempt and --model; builders commit one file per leaf, verifiers write Evidence), run under env -i. rhiza workflow check on project workflows: a leaves verifier with low_risk_models [haiku] and risky_paths [docs/**] is valid, as is a checks step's verifier with both; refused with a message naming the verifier role: low_risk_models on a leaves builder, risky_paths on a check ("set them under [step.review.verifier]"), low_risk_models on a checks fixer, low_risk_models on an agent step ("has no verifier"); refused: low_risk_models [" "] (entry 1 is blank), risky_paths [""] (entry 1 is empty), risky_paths ["a/[b"] (not a valid glob). rhiza workflow describe shows role verifier low_risk_models [haiku] and risky_paths [docs/**]; --json carries low_risk_models and risky_paths on each role (empty for the builder). go test -trimpath ./cmd/rhiza -run 'TestScripts/(workflow-describe|describe-golden)$' ok.

## Links

### Out

- verifies [verifier-tier-fields](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-fields.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
