# vm-doc-qa-1

evidence · `01M7876VK1X324DWEFBK8ZVNFS` · actor [claude-code-fcd1aaeb](../actor/claude-code-fcd1aaeb.md)

## result

PASS

## statement

At 20ae0b603 ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestVMDocLessons|TestAgentDocsConsistent': PASS. agents/README.md lists vm.md. Checked on the VM: rhiza-root.service is active, ~/code/rhiza.toml exists, /etc/apt/preferences.d/rex and /etc/apt/apt.conf.d/51unattended-no-rex exist, /usr/bin/rex belongs to the rex package, ~/orchestrator/stray-logs exists, and the built binary's daemon set help accepts maxRuns.

## Links

### Out

- verifies [vm-doc](../capability/rhiza/integrity/own-agent-docs/vm-doc.md)

## History

- `01M76NXJ53FB2GVBQPT4F4DW3R` agents/developing.md, agents/run-agents.md and agents/vm.md: the lessons for developing Rhiza, for run agents working on Rhiza, and for operating the VM
