# web-roadmap-milestone-page-verified

evidence · `01M7BM78B2ZF97P7JG4ZWC3HZR` · actor [claude-code-f2b83e2f](../actor/claude-code-f2b83e2f.md)

## result

PASS

## statement

At 25461da19 (leaf commits e16e75ca0 and 25461da19) re-ran every Checks command: go test -trimpath ./internal/web -run TestRoadmapMilestones (PASS), go test -trimpath -count=1 ./internal/web (ok), rhiza ship generate and rhiza render (exit 0). Read the diff (internal/web/roadmap.go, roadmap_test.go, roadmap.html, rhiza.css, .rhiza log/rendered only): TestRoadmapMilestones serves a two-project root and asserts no row of the other project, the finished Delta milestone only in the one in-page fold with no link, the open milestone's state word, groups in order (Building now with the run's change set, Accepted waiting, Proposed, Landed, Drafts), title then intent sub-line, bar aria-label and legend 7 landed, 2 accepted, 1 proposed, 1 withdrawn, done points in a closed fold (1 of 2 met), the landed fold '2 more landed, 1 withdrawn', and the backlog then the umbrella after the milestone and before the finished fold, which closes attempt 1's finding. No stubs or TODOs; fakes only.

## Links

### Out

- verifies [web-roadmap-milestone-page](../capability/rhiza/web-monitor/web-roadmap-milestone-page.md)

## History

- `01M79NMNGDCNJPAWRP791PRF7G` Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs
