# review-three-checks-verified

evidence · `01M7BHW6KD5SKTBE3NCKA84VA5` · actor [claude-code-492defbe](../actor/claude-code-492defbe.md)

## result

PASS

## statement

At 2bf858f08 re-ran the Checks: block: go test -trimpath ./internal/workflow -run TestAgentStepChecks (PASS: plan,code,qa parsed on an agent step; no-name, duplicate and prompt-bearing checks refused), go test -trimpath ./internal/workflow (ok), go test -trimpath ./internal/api -run TestRunFindingsCheckNames (PASS: findings and verdict with --check plan recorded as round 1 of plan in an agent step declaring plan, code, qa; --check scope refused for both ops with a message naming plan, code, qa), go test -trimpath ./cmd/rhiza -run 'TestScripts/(review|workflow)' (ok), make self-check (exit 0, only pre-existing unclaimed warnings), rhiza ship generate and rhiza render (exit 0, tree clean). Diff read: only the four named Go files plus the builder's actor log and rendered files; reviewStep accepts agent steps with declared checks via step.Checks, no hard-coded check names, no stubs or TODOs.

## Links

### Out

- verifies [review-three-checks](../capability/rhiza/workflows/built-in-workflows/review-three-checks.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
