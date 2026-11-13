# verifier-tier-defaults-qa-2

evidence · `01M79J29V6MPHQ7QSNN6G407XS` · actor [claude-code-c2ed16b5](../actor/claude-code-c2ed16b5.md)

## result

PASS

## statement

QA at 1a07d6cae with a scratch build (go build -trimpath), in a scratch Rhiza root and repo with fake rex and a fake claude logging role, cycle, leaf, attempt and --model (builders and fixers commit one file per leaf: plain -> internal/x/plain.go, risky -> .github/workflows/ci.yml; checks record majors, critics confirm them, verifiers write Evidence), run under env -i. rhiza workflow describe review prints role verifier models [opus], low_risk_models [sonnet] and the twelve risky_paths globs; the built-in build workflow verified the plain leaf on sonnet and the .github leaf on opus and build-check passed. go test -trimpath ./cmd/rhiza -run 'TestScripts/(build-workflow|workflow-describe|describe-golden)$' ok (cycle 1's build-workflow failure is fixed); go test -trimpath ./internal/workflow/... ok.

## Links

### Out

- verifies [verifier-tier-defaults](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-defaults.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
