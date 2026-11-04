# workflow-ls-verified

evidence · `01M6KH85HV5SDTT280G2J36EEK` · actor [claude-code-27a34381](../actor/claude-code-27a34381.md)

## result

PASS

## statement

Verifier claude-code-27a34381 at e1de2772 (attempt 2). Ran make test: cmd/rhiza scripts pass, including testdata/script/workflow-ls.txtar. The script sets up a personal file (XDG_CONFIG_HOME), a repo file and the built-ins, and checks that `build personal <xdg path>` is followed by `shadows repo <repo path>` and `shadows builtin builtin:build`, that mine/plan/solo/ship-change show their own layers, and that `workflow ls --json` is cmp-identical to `api call workflow.ls`. Repeated by hand with a scratch binary in a scratch repo with personal and repo build.toml, repo-only mine.toml and personal-only solo.toml: each name appears once (uniq -d empty), the winning layer and path are right, the shadowed layers are listed, and the --json and api call outputs are byte-identical. The attempt-1 finding is fixed: e1de2772 drops the rendered actor pages that came from untracked records and adds the builder's tracked actor record. `make fmt-check vet self-check`, with RHIZA_ROOT/CHANGE/RUN unset, passes (rhiza check ok). One unrelated failure: TestReadBudget (internal/api) ran roadmap/ready what-if at 61-63ms against a 50ms budget on a 4-CPU host at load average 8. It reads a generated model (gen.Default), not the repo's .rhiza, and e1de2772 changes only .rhiza files, so this is machine contention, not this leaf.

## Links

### Out

- verifies [workflow-ls](../capability/rhiza/workflows/workflow-authoring/workflow-ls.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
