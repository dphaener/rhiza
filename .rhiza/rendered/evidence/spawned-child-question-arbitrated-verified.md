# spawned-child-question-arbitrated-verified

evidence · `01M7M82V9A4TF24FP0ZK218BS1` · actor [claude-code-46715437](../actor/claude-code-46715437.md)

## result

PASS

## statement

Detached checkout at 289c5fa77, under the MemoryMax/TasksMax cap: go test -p 2 -trimpath ./internal/run -run 'TestSpawnedChildQuestion|TestConductorQuestion|TestSpawnResume' -count=1 -v exit 0 (TestSpawnedChildQuestionArbitrated included=false/true PASS); go test -p 2 -trimpath ./internal/run/... exit 0; go vet ./internal/run exit 0; rhiza ship generate exit 0 (0 files changed); rhiza render exit 0, tree clean. Diff read: pollAgent calls the new askChildren for step.Conducts() steps, which polls live spawns of the step via AgentStatus only, journals a block and reuses ask (AddQuestion, asked entry keyed on the child session, arbiter, Deliver to the child, answered entry, unblock) with no advisor or notice since arbitratesIn holds; non-question blocks are not journalled or arbitrated; reconcileSpawns untouched; a failed verdict ends the step with an error naming the question. The test asserts asked/arbiter/deliver/answered order, no StartAdvisor and no Notify in a non-autonomous run, the step finishing on the conductor's pass, the failing-verdict error detail, a spawn-ended child never polled, and a permission block not arbitrated. Only engine.go, engine_test.go and the leaf's .rhiza log/render changed; no TODOs.

## Links

### Out

- verifies [spawned-child-question-arbitrated](../capability/rhiza/workflows/arbiter-sessions/spawned-child-question-arbitrated.md)

## History

- `01M7M51W2S30H06K0BFZ6MK7PF` Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand
