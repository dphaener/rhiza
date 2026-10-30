# roadmap-app-detail-verified

evidence · `01M674MDXSF1G0Z3A2G43K8CW6` · actor [claude-code-b4ba41b7](../actor/claude-code-b4ba41b7.md)

## result

PASS

## statement

At commit 22551fb, go test ./internal/tui/roadmapapp -count=1 -v passes TestEnterOpensDetailAndEscReturns: it builds the app over fake change.show, show (question), ready and run.status results, presses j then enter on change set b and asserts the detail lists its intent, kind intention, state, leaf-one's statement and acceptance, the blocking question's statement (resolved from the ready reason 'blocked by question which-way on leaf-one') and the run's workflow ship, state stopped and step build; esc returns to the list (Backlog heading shown) with Cursor() still on b. Read the diff (internal/tui/roadmapapp/detail.go, roadmapapp.go, cmd/rhiza/tui.go): the reason regexp matches the format internal/api/plan.go emits, and leafViews reads node.create ops whose shape matches real rhiza change show --json output on this repo. make check exits 0.

## Links

### Out

- verifies [roadmap-app-detail](../capability/rhiza/roadmap-app-detail.md)

## History

- `01M673223RSVC166NJVVA13X75` The roadmap app is a cursor list of milestones and change sets that the owner acts on: open a change set's detail, plan it, ship it, accept it, answer its blocking questions, or jump to its run's Rex tab
