# workflow-guide-topic-verified

evidence · `01M6KJN0WVMQGFSZYNASSFHH5F` · actor [claude-code-6e04df56](../actor/claude-code-6e04df56.md)

## result

PASS

## statement

At cf16bd57, built rhiza and ran `rhiza guide ls`: it lists `workflows  topic  623 tokens  Writing a custom workflow...`. go test -count=1 ./internal/guide ./internal/api passed: TestWorkflowsText asserts the topic names, as list items, all eight step kinds the workflow package defines (KindAgent, Wait, Command, Arbiter, Leaves, Checks, Signoff, Workflow — the full const block in internal/workflow/workflow.go), plus layers, done/failed targets, includes and the ls/new/describe/check/run start --dry-run loop; TestGuideCommandsExist checks every `rhiza ...` reference in every guide text, workflows included, resolves to a registered operation (pendingCommands is empty). TestScripts/guide passed with guide ls counting 14 texts. Read the diff: topic text, registry entry, rhiza.md topic list, tests only.

## Links

### Out

- verifies [workflow-guide-topic](../capability/rhiza/workflows/workflow-authoring/workflow-guide-topic.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
