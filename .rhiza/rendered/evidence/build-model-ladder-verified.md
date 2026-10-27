# build-model-ladder-verified

evidence · `01M5ZMT4WX7YDB26W60EHWT184` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

Adapter clause: TestCommandModel finds --model opus before the prompt in the argv when the launch names opus, and no --model when it names none; testscript run-agent-step finds ' --model opus ' for a step with models = ["opus"] and no --model for a step with none. Engine clause: TestAgentModelLadder, with models = ["sonnet", "opus"], finds visits 1, 2 and 3 started on sonnet, opus, opus; TestLeavesModel and TestLeavesStepRetries find a leaves builder's attempts on its ladder (sonnet, opus, opus); TestLeavesRoleModel/TestStepModel check the last model is kept past the list. All passed.

## tests

go test -count=1 ./internal/adapter/claudecode -run TestCommandModel; ./internal/run -run 'TestAgentModelLadder|TestAgentModelDefault|TestLeavesModel'; ./internal/workflow -run 'TestLeavesRoleModel|TestStepModel'; ./cmd/rhiza -run 'TestScripts/run-agent-step$' (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-model-ladder](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-model-ladder.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
