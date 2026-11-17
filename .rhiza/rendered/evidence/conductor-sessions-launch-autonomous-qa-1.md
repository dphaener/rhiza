# conductor-sessions-launch-autonomous-qa-1

evidence · `01M7M99R5ZBVY1SH50X5MG6HBG` · actor [claude-code-2166ab54](../actor/claude-code-2166ab54.md)

## result

PASS

## statement

Built ./cmd/rhiza (HEAD 289c5fa77) to a scratch path; ran it in a scratch Rhiza root with a fake claude (real hook path) and fake rex, runs not started --autonomous. Workflow conduct (agent step with [[check]] qa): the conductor session launched with RHIZA_AUTONOMOUS=1 and its PreToolUse AskUserQuestion got the deny reply with the arbiter reason; a builder child from rhiza run spawn in that step also had RHIZA_AUTONOMOUS=1 and was denied. Failure case: workflow plain (agent step, no checks): the session and its spawned child had no RHIZA_AUTONOMOUS, the hook printed nothing, and the step blocked on the owner. Scoped tests: go test ./internal/api -run 'TestRunSpawn|TestSessionHook|TestStepRunnerConductorAutonomous' ok.

## Links

### Out

- verifies [conductor-sessions-launch-autonomous](../capability/rhiza/workflows/arbiter-sessions/conductor-sessions-launch-autonomous.md)

## History

- `01M7M51W2S30H06K0BFZ6MK7PF` Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand
