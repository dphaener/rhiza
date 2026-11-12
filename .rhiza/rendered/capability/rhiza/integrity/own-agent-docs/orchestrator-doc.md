# orchestrator-doc

capability · `01M76P0AFBXXQRCAEPZSZD95G9` · parent [own-agent-docs](../own-agent-docs.md)

## acceptance

A test in cmd/rhiza reads agents/orchestrator.md and finds each of those rules (git fetch, pkill, uncommitted on main, maxRuns, merge origin/main, --squash, quote the owner, milestone planning with the owner); it and the agent-docs consistency test pass with go test -trimpath ./cmd/rhiza -run on their names

## statement

agents/orchestrator.md, listed in agents/README.md, tells the orchestrator agent how to run Rhiza's own development: no git fetch in the main clone from the watch loop and no pkill by pattern, the main tree stays clean and nothing is left uncommitted on main (every model write on its own branch and PR), lower rhiza.daemon.maxRuns when the VM is loaded, merge origin/main into a conflicting PR and never rebase and force-push, merge by hand with gh pr merge --squash, quote the owner exactly, and plan a milestone only with the owner

## Links

### In

- verifies [orchestrator-doc-verified](../../../../evidence/orchestrator-doc-verified.md)
- verifies [orchestrator-doc-qa-1](../../../../evidence/orchestrator-doc-qa-1.md)
- verifies [pr-496-orchestrator-agent-doc](../../../../evidence/pr-496-orchestrator-agent-doc.md)

## History

- `01M76NX8F3G4KXGHQ9G93TEH4Z` agents/orchestrator.md: how the orchestrator runs Rhiza's own development, consulting the owner before judgment calls and acting autonomously only when granted
