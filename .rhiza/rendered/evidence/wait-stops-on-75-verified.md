# wait-stops-on-75-verified

evidence · `01M6AH99H0609F3WM6278B3VP2` · actor [claude-code-9e8c105a](../actor/claude-code-9e8c105a.md)

## result

PASS

## statement

At b19150bf: go test ./internal/run -run 'TestWait|TestCommandStopExit' -count=1 -v passes, including the new TestWaitStopExitStopsRun. In that test the condition reports 1, then 75 with detail 'rebase conflict', then 0. The run stops with Stopped='rebase conflict', Ended=false, Current='ci' and no step outcome. The first exit 1 reads as not held, so the condition was checked twice. A second e.Run resumes, checks the wait a third time, gets exit 0 and succeeds with ci:pass. In the diff, engine.wait uses a new Conditioner interface and returns result{stopped: detail} on StopExit without ending the step; any other code is held only when it is 0. api stepRunner.Condition runs step.Until with stderr merged into stdout and takes the detail from lastStopLine, the same helper Command now uses for its last 'stop: ' line. go test ./... passes in every package; on its first run make check hit the read-timing budget flake, and on a rerun fmt-check, vet and self-check are ok.

## Links

### Out

- verifies [wait-stops-on-75](../capability/rhiza/workflows/workflow-definitions/wait-step/wait-stops-on-75.md)

## History

- `01M6A0JJDJ0FXD58K4WV88FV0A` A ship-change wait settles a PR that conflicts after main moves by rebasing it itself, and stops for the owner on a code conflict
