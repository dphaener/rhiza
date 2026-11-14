# rhiza-repo-conductor-default-qa-2

evidence · `01M7DTKR7R3PE71HDQ5F2H41KX` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: In this repo, workflow describe ship-change resolves to builtin:ship-change; .rhiza/workflows/ship-change.toml is gone. describe ship-rhiza shows the env GOTMPDIR, SHIP_HEAD_RHIZA and SHIP_HEAD_BUILD, go-cleanup (sh .rhiza/workflows/go-cleanup.sh), the 10 included steps and handback:ship-change/conduct. bugfix shows the same. workflow check ship-rhiza, bugfix and plan-bug are valid, with only the bypass warnings AGENTS.md expects. make self-check in a clean detached worktree at HEAD exits 0. TestShipFixPromptsOneLine passes.

## Links

### Out

- verifies [rhiza-repo-conductor-default](../capability/rhiza/workflows/workflow-overrides/rhiza-repo-conductor-default.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
