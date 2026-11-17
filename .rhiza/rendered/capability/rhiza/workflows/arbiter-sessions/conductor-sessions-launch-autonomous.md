# conductor-sessions-launch-autonomous

capability · `01M7M5E1000DMWWQVN057JWBV5` · parent [arbiter-sessions](../arbiter-sessions.md)

## acceptance

go test -trimpath ./internal/api -run 'TestRunSpawn|TestSessionHook|TestStepRunnerConductorAutonomous': in a run not started --autonomous whose current step is an agent step declaring checks, a new TestStepRunnerConductorAutonomous shows the session stepRunner.StartAgent starts for that step recorded Autonomous with RHIZA_AUTONOMOUS=1 in its launch env, and a TestRunSpawn case shows a child spawned in that step recorded Autonomous too; a hook PreToolUse AskUserQuestion payload from either is answered with the deny reply carrying autonomousAskReason and sets the record's asked mark; (failure case) in the same non-autonomous run a session started for an agent step without checks, and a child spawned in such a step, are not autonomous and their AskUserQuestion is not denied.

## files

internal/api/run.go (stepRunner.StartAgent), internal/api/run_spawn.go (runSpawn's SessionStartRequest), internal/api/run_spawn_test.go, internal/api/session_ask_test.go

## must_not

Change the run's State.Autonomous, daemonStart, the signoff step or the gates; mark sessions of non-conductor steps autonomous; change autonomousAskReason or the hook's deny format

## statement

The session of a conductor step (an agent step that declares checks), and every child session rhiza run spawn starts while the run is in such a step, is recorded autonomous whether or not the run started --autonomous: its launch carries RHIZA_AUTONOMOUS=1, its AskUserQuestion is denied with the arbiter reason and kept as the asked mark, and no Rex blocked notice is written for its question.

## Links

### Out

- decided-by [conductor-questions-to-arbiter](../../../../decision/conductor-questions-to-arbiter.md)

### In

- depends-on [spawned-child-question-arbitrated](spawned-child-question-arbitrated.md)
- verifies [conductor-sessions-launch-autonomous-verified](../../../../evidence/conductor-sessions-launch-autonomous-verified.md)
- verifies [conductor-sessions-launch-autonomous-qa-1](../../../../evidence/conductor-sessions-launch-autonomous-qa-1.md)

## History

- `01M7M51W2S30H06K0BFZ6MK7PF` Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand
