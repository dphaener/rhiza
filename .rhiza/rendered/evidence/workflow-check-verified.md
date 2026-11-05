# workflow-check-verified

evidence · `01M6KEP69GE6RENP0XWT1S5Z4Q` · actor [claude-code-5ae350b7](../actor/claude-code-5ae350b7.md)

## result

PASS

## statement

Verifier claude-code-5ae350b7 at 63f55027. go test ./cmd/rhiza -run TestScripts/(workflow-check|describe-golden)$ passes: workflow-check.txtar runs rhiza workflow check on every built-in (build, plan, review, ship-change, each exit 0 'is valid') and on the internal/workflow testdata files (embedded copies are byte-identical to internal/workflow/testdata/*.toml: checked with diff). valid.toml exits 0; bad-outcome, missing-target and unknown-kind exit non-zero with 'step <name>: ...' naming the problem; a missing include is reported with its including step and is expanded once the repo layer has it; several problems are all reported; --json is cmp-equal to rhiza api call workflow.check for a bad and a good workflow. Repeated by hand with a scratch build in a scratch repo: same verdicts, exit 0/1, --json bytes cmp-equal to api call. Read the diff: Validate reports every per-step problem, then defers to ParseWith with the run-start resolver for includes. make test (RHIZA_ROOT etc. unset): every package ok; only the timing budgets (TestReadBudget, slice TestBudget, TestGeneratedModel, none in code this leaf touches) missed on a shared VM at load average 5-12 on 4 cores. make self-check fails on stale rendered actor files (043bcc3c, 1cc8d4e7, 1205bb4f) whose records are untracked: a render problem that predates this leaf (it also affects earlier leaves) and is already fixed after this commit by e1de2772 'Render only the actors the committed log records'.

## Links

### Out

- verifies [workflow-check](../capability/rhiza/workflows/workflow-authoring/workflow-check.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
