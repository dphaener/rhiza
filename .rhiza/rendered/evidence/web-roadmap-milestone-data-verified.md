# web-roadmap-milestone-data-verified

evidence · `01M7BKA84HPKKXRRTQGW5DRC35` · actor [claude-code-c48c1aee](../actor/claude-code-c48c1aee.md)

## result

PASS

## statement

At 5ec8ca90a re-ran every Checks: command, all exit 0: go test -trimpath ./internal/api -run TestMilestonesOperation (PASS), go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden (ok; re-running with -update left the golden unchanged), rhiza ship generate (0 files changed), rhiza render (no drift), make self-check (exit 0, only pre-existing unclaimed warnings). Read the diff: internal/api/milestones.go adds a separate read-only milestones operation; roadmap is untouched. The test fixture covers each acceptance point: milestone A has draft, proposed, accepted, landed and withdrawn change sets, all with their states; A1 (landed plus withdrawn) is met; A2 (accepted) is not met; A is building, B planned, C designing, D finished; finished=1; roadmap still lists only the proposed and accepted change sets. The cli.md docs and the describe golden were regenerated, there are no TODOs, and the only other files are the actor record and the rendered files. A scratch build of rhiza milestones --json on this repo returned 36 milestones with 6 finished.

## Links

### Out

- verifies [web-roadmap-milestone-data](../capability/rhiza/web-monitor/web-roadmap-milestone-data.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
