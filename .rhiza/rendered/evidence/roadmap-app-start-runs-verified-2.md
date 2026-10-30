# roadmap-app-start-runs-verified-2

evidence · `01M676T8ETCNEXG7ER1E7XFAHA` · actor [claude-code-4c775004](../actor/claude-code-4c775004.md)

## result

PASS

## statement

At 8fce720 (attempt 2, fix for review finding f1) ran go test ./internal/tui/roadmapapp -count=1: TestPlanRunStartsInteractiveInNewTab (p,y gives one rex window new call running rhiza run start <cs> --workflow plan with no --autonomous, view says which tab opened, app stays on the list), TestShipRunIsAutonomousWithAutoMode (s,y gives run start <cs> --autonomous --permission-mode auto), TestNDeclinesAndCallsNothing (n calls nothing) and TestOFocusesLiveRunOrSaysNone (o focuses the live run's tab or says No live run) all PASS with a fake exec, no real Rex. go test ./internal/api -run RunOpen PASS: run.open now resolves the change ref through resolveChangeSet, so a lowercase unique prefix finds the live run, the response carries the full ID, and the rex calls are block inspect then focus. Read the diff: only run_open.go and its test changed. go test ./..., fmt-check and vet pass; make self-check fails only on stale .rhiza/rendered actor/index files from later review record commits, not this leaf's code.

## Links

### Out

- verifies [roadmap-app-start-runs](../capability/rhiza/roadmap-app-start-runs.md)

## History

- `01M673223RSVC166NJVVA13X75` The roadmap app is a cursor list of milestones and change sets that the owner acts on: open a change set's detail, plan it, ship it, accept it, answer its blocking questions, or jump to its run's Rex tab
