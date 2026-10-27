# run-program-status-verified

evidence · `01M5ZMT8NRXD1J1JYV0VV7C78M` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

TestProgramStatusBuild captures the OSC 7501 bytes a run writes to a fake Rex buffer through a build of three leaves and finds, in order: root working with 'step build' first; per leaf a child record leaf/N with builder working then verifier, then clear; root progress from 0 rising to 100; a blocked child for leaf/2 while its fake verifier waits on a question; root done last. TestProgramStatusError finds a failing leaf run and a failing command run end with root error. All passed.

## tests

go test -count=1 ./internal/api -run 'TestProgramStatusBuild|TestProgramStatusError|TestRexRecord'; ./internal/run -run 'TestEngineStatus|TestDeriveStatus' (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [run-program-status](../capability/rhiza/workflows/workflow-runs/run-program-status.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
