# vm-doc

capability · `01M76P0B97WK63HSDNR4DMB2BT` · parent [own-agent-docs](../own-agent-docs.md)

## acceptance

A test in cmd/rhiza reads agents/vm.md and finds each lesson (rhiza-root.service, maxRuns, the Rex pin, stray-logs, inject delay); it and the agent-docs consistency test pass with go test -trimpath ./cmd/rhiza -run on their names

## statement

agents/vm.md, listed in agents/README.md, holds the lessons for operating the VM that runs Rhiza's daemon: one root ~/code under rhiza-root.service, lower rhiza.daemon.maxRuns when load is high, the Rex package is pinned because Ubuntu's Perl rex package replaced it once, worktree resets move stray .rhiza/log entries to ~/orchestrator/stray-logs/, and flakes are reproduced by injecting delay, never by CPU load on the shared host

## Links

### In

- verifies [vm-doc-verified](../../../../evidence/vm-doc-verified.md)
- verifies [vm-doc-qa-1](../../../../evidence/vm-doc-qa-1.md)
- verifies [pr-497-role-docs-for-rhiza-agents](../../../../evidence/pr-497-role-docs-for-rhiza-agents.md)

## History

- `01M76NXJ53FB2GVBQPT4F4DW3R` agents/developing.md, agents/run-agents.md and agents/vm.md: the lessons for developing Rhiza, for run agents working on Rhiza, and for operating the VM
