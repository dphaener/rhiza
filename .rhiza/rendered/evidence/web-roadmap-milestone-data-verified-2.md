# web-roadmap-milestone-data-verified-2

evidence · `01M7BTMEGPDF3K0GH4B887J7D2` · actor [claude-code-4946b474](../actor/claude-code-4946b474.md)

## result

PASS

## statement

Verified leaf commits 5ec8ca90a and c1e692b2a (attempt 2). Re-ran the Checks block: go test -trimpath ./internal/api -run 'TestEveryResponseIsHumaner|TestReadOnlyMarker' exit 0; go test -trimpath -race ./internal/api -run TestMilestones -count=3 exit 0 (3 PASS); go test -trimpath ./internal/web/... exit 0. The whole go test -trimpath ./internal/api hit the default 10m timeout twice on a host at load ~10 (first in TestRootIdlePollBudget waiting for a quiet host, then with budgets skipped); with -skip Budget -timeout 25m it passes: 713 tests, 849s, the slowest TestLinkEndpointKinds 202s, none of them touched by this leaf. Acceptance: TestMilestonesOperation passes and its fixture asserts A building with A1 met, A2 not and five change sets with states draft/proposed/accepted/landed/withdrawn; B planned, C designing, D finished, finished count 1; roadmap still only A proposed, A accepted and B proposed. go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden exit 0 with milestones in the golden. Diff read: milestones.go (operation, Human output, a mutex around EachInFlight's map), milestones_test.go, read_only_test.go, describe golden, docs/reference/cli.md, .rhiza records; nothing else, no TODOs or stubs, fakes only.

## Links

### Out

- verifies [web-roadmap-milestone-data](../capability/rhiza/web-monitor/web-roadmap-milestone-data.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
