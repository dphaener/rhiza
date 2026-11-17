# workflow-keys-topic-verified

evidence · `01M7MM4SK3890BB5FC673V4NYB` · actor [claude-code-f844299d](../actor/claude-code-f844299d.md)

## result

PASS

## statement

At 7332316a9 re-ran every Checks: command, each exit 0: go test -p 2 -trimpath ./internal/guide -run 'TestWorkflowKeys|TestWorkflowsText|TestConduct|TestList' -v (TestWorkflowKeys, TestConductReadsWorkflowKeys, TestWorkflowKeysCatchesDrift, TestWorkflowsText, TestList, TestListOrder pass); go test ./internal/guide/...; go test ./cmd/rhiza -run TestScripts/guide (guide ls count 40, workflow-keys listed); rhiza ship generate; rhiza render. Read the diff: TestWorkflowKeys iterates workflow.FileKeys, SharedKeys, RoleKeys() and KindDocs() (keys, tables, outcomes) and checks TopicCap, the defaults, the conductor-step line and the rhiza workflow new phrases; the failure case is covered by TestWorkflowKeysCatchesDrift (added timezone/retries/effort/lint and a removed gate and stall are each reported by name). Workflows topic names gate and rhiza guide workflow-keys under TopicCap; conduct Read first names rhiza guide workflow-keys under TopicCap. Spot-checked the topic's defaults against internal/workflow/spec.go and the skeleton against internal/api/workflow_new.go: they match. Only guide text, its tests and guide.txtar changed; internal/workflow untouched; no TODOs.

## Links

### Out

- verifies [workflow-keys-topic](../capability/rhiza/guide/agent-guide/workflow-keys-topic.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
