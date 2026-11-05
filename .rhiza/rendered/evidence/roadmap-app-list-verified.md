# roadmap-app-list-verified

evidence · `01M6747M3NR3XG276B90Z5AN4B` · actor [claude-code-9183fc13](../actor/claude-code-9183fc13.md)

## result

PASS

## statement

At 4444bfc, go test ./internal/tui/roadmapapp -count=1 -v passed: TestViewIsAListNotJSON builds the app over fake roadmap and ready results and asserts the milestone heading (rhiza / milestone-5), the Backlog heading, rows with short ID, [state] and title, the not-ready reasons, and no { or } in the view; TestCursorSkipsHeadings asserts j moves from the first row to the second and across the Backlog heading to the third, stops at the end, and k moves back; TestReloadKeepsCursorAndResultIsRoadmapBytes covers r. go test ./cmd/rhiza -run TestScripts/tui passed, including tui-read-apps.txtar, where rhiza tui app roadmap --once compares equal (cmp) to rhiza roadmap --json. A scratch build on this repo's model also gave byte-identical output for tui app roadmap --once and roadmap --json. Reading the diff: roadmap moved from readapps into internal/tui/roadmapapp, wired in cmd/rhiza/tui.go, and q/esc quit. fmt-check, vet and go test ./... pass, with the slice and modelapp timing tests passing when rerun alone. make self-check fails only on rendered files for actor records this run wrote after the commit, not on the leaf.

## Links

### Out

- verifies roadmap-app-list (retired)

## History

- `01M673223RSVC166NJVVA13X75` The roadmap app is a cursor list of milestones and change sets that the owner acts on: open a change set's detail, plan it, ship it, accept it, answer its blocking questions, or jump to its run's Rex tab
