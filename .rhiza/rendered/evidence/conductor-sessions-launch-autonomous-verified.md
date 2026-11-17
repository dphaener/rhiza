# conductor-sessions-launch-autonomous-verified

evidence · `01M7M78JJH0DS8CJYWDG18881J` · actor [claude-code-176f7c67](../actor/claude-code-176f7c67.md)

## result

PASS

## statement

Detached checkout at faeaeef60, all under the capped scope: acceptance go test -p 2 -trimpath ./internal/api -run 'TestRunSpawn|TestSessionHook|TestStepRunnerConductorAutonomous' exit 0 (TestStepRunnerConductorAutonomous/Conductor+Plain and TestRunSpawn/ConductorStepAutonomous+PlainStepNotAutonomous ran and passed); go vet ./internal/api ./internal/workflow exit 0; go test ./internal/workflow/... exit 0; go test ./internal/api/... exit 0 (307s); rhiza ship generate and rhiza render exit 0 with no drift. Diff read: StartAgent ORs step.Conducts() into req.autonomous; runSpawn ORs spawnInConductor (def.Steps lookup of the spawn step, the same as reviewStep) into the child's request; workflow.Step.Conducts = agent step with checks. Tests assert the record's Autonomous, RHIZA_AUTONOMOUS=1 in the launch env, the hook's deny reply with autonomousAskReason and the asked mark, and for plain agent steps no autonomy, no output and no asked mark. State.Autonomous, daemonStart, signoff, gates, autonomousAskReason and the deny format are untouched. No TODOs. The one file outside the leaf's list is internal/workflow/workflow.go (+6 lines, the shared Conducts predicate the run notes ask for).

## Links

### Out

- verifies [conductor-sessions-launch-autonomous](../capability/rhiza/workflows/arbiter-sessions/conductor-sessions-launch-autonomous.md)

## History

- `01M7M51W2S30H06K0BFZ6MK7PF` Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand
