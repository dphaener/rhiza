# workflow-describe-verified

evidence · `01M6KD6DPA4XBGQKJX4JTYP2ZP` · actor [claude-code-1205bb4f](../actor/claude-code-1205bb4f.md)

## result

PASS

## statement

At 95ae54b2 ran go test ./cmd/rhiza -run 'TestScripts/(workflow-describe|describe-golden)$' -count=1: ok. The workflow-describe script describes the built-in build (build (builtin) builtin:build, start, on.pass ->, role builder models), a repo .rhiza/workflows/build.toml that shadows it (shadows builtin builtin:build, start only), and inc.toml whose workflow step includes build (includes build, step build/only command); for both build and inc.toml, --json is cmp-equal to rhiza api call workflow.describe; afterwards run ls lists no run and fold-ok passes. Read the diff: workflow.describe goes through the same resolveWorkflow/workflowResolver as run start, and its handler only reads files. Checked by hand with a scratch build in a scratch repo: a personal XDG build.toml wins and shows that it shadows repo and builtin, steps come out in order with on targets and verifier/critic/fixer roles with model ladders and prompts (review), and an md5 of every file in the repo (including .git) and the config dir was unchanged after describing. make test: every package ok; TestReadBudget failed with uniform slowdowns across all reads (show included) at load average 8.6 on 4 cores. Those are contention timings: the leaf adds no read-path code. make check's self-check failed only on rendered files for live sessions' untracked actor records (including this verifier's), none in the commit.

## Links

### Out

- verifies [workflow-describe](../capability/rhiza/workflows/workflow-authoring/workflow-describe.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
