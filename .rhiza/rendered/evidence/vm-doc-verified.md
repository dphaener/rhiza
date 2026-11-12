# vm-doc-verified

evidence · `01M78703RMT8XFEF1MS9PERZM6` · actor [claude-code-51f4313b](../actor/claude-code-51f4313b.md)

## result

PASS

## statement

At 278aebb85 (detached checkout) ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestVMDocLessons|TestAgentDocsConsistent' -v: TestVMDocLessons (new cmd/rhiza/vmdoc_test.go reads agents/vm.md and checks rhiza-root.service, rhiza.daemon.maxRuns, 'The Rex package is pinned', stray-logs, 'injecting a delay'), TestAgentDocsConsistent and its fixtures all PASS. Read the diff: agents/vm.md covers one root ~/code under rhiza-root.service, lowering rhiza.daemon.maxRuns on high 5-minute load, the Rex apt pin after Ubuntu's Perl rex replaced it, worktree resets moving stray .rhiza/log entries to ~/orchestrator/stray-logs/, and reproducing flakes by injecting delay not CPU load; agents/README.md indexes vm.md. Cross-checked on the VM: rhiza-root.service is active and /etc/apt/preferences.d/rex and /etc/apt/apt.conf.d/51unattended-no-rex exist. Commit stages only the builder's own actor record.

## Links

### Out

- verifies [vm-doc](../capability/rhiza/integrity/own-agent-docs/vm-doc.md)

## History

- `01M76NXJ53FB2GVBQPT4F4DW3R` agents/developing.md, agents/run-agents.md and agents/vm.md: the lessons for developing Rhiza, for run agents working on Rhiza, and for operating the VM
