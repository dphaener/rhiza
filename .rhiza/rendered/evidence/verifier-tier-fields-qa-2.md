# verifier-tier-fields-qa-2

evidence · `01M79J29WGRGYT85B5WPARQT7S` · actor [claude-code-c2ed16b5](../actor/claude-code-c2ed16b5.md)

## result

PASS

## statement

QA at 1a07d6cae with a scratch build (go build -trimpath), in a scratch Rhiza root and repo with fake rex and a fake claude logging role, cycle, leaf, attempt and --model (builders and fixers commit one file per leaf: plain -> internal/x/plain.go, risky -> .github/workflows/ci.yml; checks record majors, critics confirm them, verifiers write Evidence), run under env -i. A project review.toml setting low_risk_models and risky_paths on a checks step's verifier loaded and drove the run, and rhiza workflow describe review shows both fields on the verifier role. go test -trimpath ./internal/workflow/... and ./cmd/rhiza TestScripts/(workflow-describe|describe-golden) ok.

## Links

### Out

- verifies [verifier-tier-fields](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-fields.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
