# workflow-overrides-verified

evidence · `01M6K6J2QZK31AHETZY0ZF37DA` · actor [claude-code-5f6c9898](../actor/claude-code-5f6c9898.md)

## result

PASS

## statement

At a40e9c23 I read the diffs of c94a9094 (lookupWorkflow and personalWorkflowsDir in internal/api/run.go, shared by resolveWorkflow for run start and for workflow steps) and a40e9c23 (cmd/rhiza/testdata/script/workflow-personal-override.txtar). I ran go test ./cmd/rhiza -run TestScripts/workflow-personal-override -count=1 (ok) and go test ./internal/api -run Workflow -count=1 (ok). The script uses a hermetic HOME, sets XDG_CONFIG_HOME empty and then to $WORK/xdg, and runs real run start --foreground. Its cases: the repo's mine.toml runs and is recorded; ~/.config/rhiza/workflows/mine.toml shadows it while XDG_CONFIG_HOME is unset; $XDG_CONFIG_HOME/rhiza/workflows/mine.toml shadows both; a workflow step in outer naming mine runs the personal step; a repo build.toml shadows the built-in build; a personal build.toml shadows the repo file and the built-in; with no file the built-in runs. Each case greps the run journal for workflow_path set to the file's path, or builtin:build for the built-in, which resolves attempt 1's finding f1. In make check, fmt-check, vet, every functional test and make self-check passed. Only the timing budgets TestReadBudget and TestBudget missed their 50ms what-if medians, under load average ~10 from other worktrees' suites, on read paths this leaf does not touch.

## Links

### Out

- verifies [workflow-overrides](../capability/rhiza/workflows/workflow-overrides.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
