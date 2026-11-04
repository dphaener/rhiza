# workflow-describe-verified-2

evidence · `01M6KQG79SAKYB0KSMYD8KEJN1` · actor [claude-code-b3491ab1](../actor/claude-code-b3491ab1.md)

## result

PASS

## statement

Verifier claude-code-b3491ab1 at e4e41daf (attempt 2, fixing quality f2). Read the diff: runWorkflowDescribe now picks the repo layer and finds lower layers by the resolved ref, not def.Name, and defaults an empty ref to defaultWorkflow as resolveWorkflow does. Ran go test ./cmd/rhiza -run TestScripts/workflow-describe$ -count=1: ok. The script describes built-in build (build (builtin) builtin:build, on.pass ->, role builder models); a repo .rhiza/workflows/build.toml override shows '(repo)' with 'shadows builtin builtin:build'; inc.toml, which includes build, shows 'includes build' and 'build/only command'; --json cmp-equals rhiza api call workflow.describe for both build and inc.toml; run ls shows no run and fold-ok passes. New cases: mine.toml with name = "custom" is reported as repo with no shadows; a personal copy shadows the repo one. Reproduced f2 by hand with a scratch build: describe mine now prints 'custom (repo) .../r/.rhiza/workflows/mine.toml', and a repo build.toml shows 'shadows builtin builtin:build'. go test ./internal/api -run Workflow: ok. make check: fmt-check, vet and test passed; self-check failed only on stale .rhiza/rendered actor files from review-run records, not this leaf. A later make test failed only the timing budgets (TestReadBudget, slice TestBudget, modelapp TestGeneratedModel) at load average 11-18; every other package was ok, and this diff touches no read path.

## Links

### Out

- verifies [workflow-describe](../capability/rhiza/workflows/workflow-authoring/workflow-describe.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
