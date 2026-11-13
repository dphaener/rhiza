# verifier-tier-choice-qa-2

evidence · `01M79J29SPPBR1VKNKKF6QAD4G` · actor [claude-code-c2ed16b5](../actor/claude-code-c2ed16b5.md)

## result

PASS

## statement

QA at 1a07d6cae with a scratch build (go build -trimpath), in a scratch Rhiza root and repo with fake rex and a fake claude logging role, cycle, leaf, attempt and --model (builders and fixers commit one file per leaf: plain -> internal/x/plain.go, risky -> .github/workflows/ci.yml; checks record majors, critics confirm them, verifiers write Evidence), run under env -i. A build run then a review run with a project review.toml (one check, verifier models [opus], low_risk_models [sonnet], risky_paths [*_test.go, .github/**]). Build: plain verifier sonnet, risky verifier opus. Review with majors on both leaves in cycle 1: cycle 2 first fix, journalled leaf-session plain verifier attempt 2 sonnet, risky verifier attempt 2 opus. With plain's first cycle-2 verifier failing: plain verifier attempt 3 (second try of cycle 2) opus. With a major on plain again in cycle 2: cycle 3 fix, plain verifier attempt 5 opus. Fixers stay on opus. Both runs succeeded. go test -trimpath ./internal/run -run 'TestVerifierTier|TestVerifierRisky' ok.

## Links

### Out

- verifies [verifier-tier-choice](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-choice.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
