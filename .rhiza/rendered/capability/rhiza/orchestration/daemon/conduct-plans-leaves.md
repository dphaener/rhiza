# conduct-plans-leaves

capability · `01M79KPD9T4QSDZCWH4NQ577EA` · parent [daemon](../daemon.md)

## acceptance

go test -trimpath ./internal/api -run 'TestPickNoLeaves|TestDaemonPicks' shows a leafless, milestone-serving change set dispatched to a workflow starting with an agent step that declares checks picked, the same change set dispatched to a workflow starting with a leaves step skipped as no-leaves, a leafless backlog change set still skipped, and an umbrella change set still not ready. go test -trimpath ./cmd/rhiza -run TestScripts/run-ship-needs-leaves shows rhiza run start --workflow ship-change still refusing a leafless change set with 'has no leaves' and next 'rhiza run start --workflow plan' when ship-change starts with a command or leaves step, and starting the run on the same leafless change set when the repo's ship-change starts with an agent step that declares checks (a fake agent, nothing reaching Claude Code, GitHub or Rex). Touches internal/api/daemon.go, internal/api/daemon_picks_test.go, internal/api/run.go, cmd/rhiza/testdata/script/run-ship-needs-leaves.txtar.

## statement

The daemon picks an accepted, ready change set that serves a milestone and has no leaves when its dispatch resolves to the conductor default (a workflow whose first step is an agent step that declares checks), instead of skipping it; backlog change sets with no leaves and change sets dispatched to a classic workflow keep the no-leaves skip, and umbrella change sets stay not ready. rhiza run start lets such a leafless change set start: its no-leaves refusal (checkShipLeaves) applies to ship-change only when the workflow's first step is not an agent step that declares checks, so a ship-change overridden by the conductor default starts on a leafless change set and a classic ship-change still refuses one with the plan-first next command.

## Links

### Out

- depends-on [review-three-checks](../../workflows/built-in-workflows/review-three-checks.md)
- decided-by [conductor-plans-leaves](../../../../decision/conductor-plans-leaves.md)
- decided-by [conductor-checks-gate-pick](../../../../decision/conductor-checks-gate-pick.md)

### In

- verifies [conduct-plans-leaves-verified](../../../../evidence/conduct-plans-leaves-verified.md)
- verifies [conduct-plans-leaves-qa-1](../../../../evidence/conduct-plans-leaves-qa-1.md)
- verifies [conduct-plans-leaves-verified-2](../../../../evidence/conduct-plans-leaves-verified-2.md)
- verifies [conduct-plans-leaves-qa-2](../../../../evidence/conduct-plans-leaves-qa-2.md)
- verifies [pr-575-a-rhiza-conduct-skill-conducts](../../../../evidence/pr-575-a-rhiza-conduct-skill-conducts.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
