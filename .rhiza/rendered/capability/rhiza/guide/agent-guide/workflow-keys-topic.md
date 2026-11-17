# workflow-keys-topic

capability · `01M7MJXM02SHB19P39VYGDNXSH` · parent [agent-guide](../agent-guide.md)

## acceptance

go test -trimpath ./internal/guide -run 'TestWorkflowKeys|TestWorkflowsText|TestConduct|TestList' shows the workflow-keys topic listed and under TopicCap; the test iterates workflow.FileKeys, workflow.SharedKeys, workflow.RoleKeys() and workflow.KindDocs() (each kind, its keys, tables and outcomes) and requires the topic to name each; it requires the defaults above, [[step.<name>.check]], rhiza workflow new with --from, handback:build, max_visits = 3 and the refusal of an existing file; TestWorkflowsText requires the workflows topic to name the gate kind and rhiza guide workflow-keys and the conduct topic's Read first names rhiza guide workflow-keys, both under TopicCap; (failure case) the test fails naming the key or kind when spec.go documents one the topic does not name; go test -trimpath ./cmd/rhiza -run TestScripts/guide passes.

## files

internal/guide/text/workflow-keys.md, internal/guide/workflow_keys_test.go, internal/guide/guide.go, internal/guide/guide_test.go, internal/guide/text/workflows.md, internal/guide/workflows_test.go, internal/guide/text/conduct.md, internal/guide/conduct_test.go, cmd/rhiza/testdata/script/guide.txtar

## must_not

Change internal/workflow, rhiza workflow new or any workflow file: guide text and its tests only; grow the conduct or workflows topic past TopicCap; drop anything the workflows topic says today

## statement

The guide has a workflow-keys topic, named in the conduct topic's Read first and in the workflows topic, listing every top-level file key, shared step key, step kind (gate included) with its keys, sub-tables and outcomes, and role sub-table key that internal/workflow/spec.go documents, with the default each takes when unset (name the file's base name, start the first step, handback_cap 3, max_visits once, stall the watchdog's default, agent and permission_mode the run's, models the agent's default with the last kept, interval 10s or 30s with until, until_when idle, split auto, an autonomy stage rhiza.autonomy.<stage>, a cleanup name cleanup-<n>, a built-in arbiter when none is declared), that an agent step declaring [[step.<name>.check]] is a conductor step, and what rhiza workflow new writes: with --from a copy of the resolved workflow with its name set, without it a skeleton starting at an agent step build whose command step test is a placeholder that fails on purpose, with max_visits 3 and on.fail handback:build; it refuses an existing file and validates as workflow check does. The workflows topic's step kinds list names gate.

## Links

### Out

- decided-by [conduct-guide-three-topics](../../../../decision/conduct-guide-three-topics.md)

### In

- verifies [workflow-keys-topic-verified](../../../../evidence/workflow-keys-topic-verified.md)
- verifies [workflow-keys-topic-qa-1](../../../../evidence/workflow-keys-topic-qa-1.md)
- verifies [workflow-keys-topic-qa-2](../../../../evidence/workflow-keys-topic-qa-2.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
