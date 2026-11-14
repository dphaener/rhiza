# conduct-skill-verified-2

evidence · `01M7BS2TDGV007PE4NT8FNXRDR` · actor [claude-code-93aa3402](../actor/claude-code-93aa3402.md)

## result

PASS

## statement

At 1c7e620ee (attempt 2, fix for qa f2): go test -trimpath ./internal/api -run 'TestRunSpawn|TestRunFindings|TestReviewCheck|TestSession|TestDescribe|TestRunReview|TestGuideInstall' exit 0; go test -trimpath ./internal/guide/... exit 0 and -run TestConduct PASS; go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-spawn|describe-golden)$' exit 0; rhiza ship generate (0 files changed) and rhiza render leave the tree clean. The whole go test ./internal/api/... timed out at the 10m package limit (TestDaemonAutoUpdateLatest 4s in), matching the Checks block's recorded timeout and unrelated to the diff. The diff adds run spawn --check and spawnCheck (check from a check-<check> role or --check, refused when the step does not declare it), wired through session.go into RHIZA_CHECK/RHIZA_CYCLE (session_run.go:561), tested by TestRunSpawnCheck (check role, critic, unknown check, no-checks step); conduct.md and TestConduct name --check and RHIZA_CHECK; describe-golden and cli.md are regenerated. No stubs or TODOs, and nothing outside the leaf's slice besides the generated docs and rendered files.

## Links

### Out

- verifies [conduct-skill](../capability/rhiza/guide/agent-guide/conduct-skill.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
