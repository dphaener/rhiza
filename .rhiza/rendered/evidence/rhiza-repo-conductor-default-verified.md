# rhiza-repo-conductor-default-verified

evidence · `01M7C870H3SHQXBFVJWD9VPE9D` · actor [claude-code-1ad58c78](../actor/claude-code-1ad58c78.md)

## result

PASS

## statement

Verifier attempt 2 at 0383f0256, with a scratch build of this commit and RHIZA_ROOT set to the checkout. rhiza workflow describe ship-change resolves to builtin:ship-change (exit 0). describe ship-rhiza (repo) shows env GOTMPDIR, SHIP_HEAD_RHIZA and SHIP_HEAD_BUILD, cleanup go-cleanup: sh .rhiza/workflows/go-cleanup.sh, the builtin's 10 steps and every hand-back route as handback:ship-change/conduct. describe bugfix shows the same go-cleanup and the same handback:ship-change/conduct routes. rhiza workflow check ship-rhiza, bugfix and plan-bug: exit 0, with only the autonomous or ready-change-set bypass warnings that AGENTS.md documents. .rhiza/workflows/ship-change.toml is gone and dispatch.toml has [ready] feature=ship-rhiza and bug=bugfix. In a clean detached worktree: rhiza check exit 0, make self-check exit 0. go test -trimpath ./internal/workflow/... ok; TestShipFixPromptsOneLine and TestIncludeSignoffHandback PASS; TestScripts/(workflow-describe|describe-golden) ok; internal/api TestWorkflow* ok. TestRequestExposedWorkflows also fails at HEAD~2, before this leaf, so it is not caused here. The extra files (include.go signoff hand-back, describe [env], AGENTS.md, docs) answer attempt 1's findings and are tested. The diff has no TODOs or stubs.

## Links

### Out

- verifies [rhiza-repo-conductor-default](../capability/rhiza/workflows/workflow-overrides/rhiza-repo-conductor-default.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
