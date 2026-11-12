# run-agents-doc-qa-1

evidence · `01M7876VG0CWGWPGYC7S9CJ9SB` · actor [claude-code-fcd1aaeb](../actor/claude-code-fcd1aaeb.md)

## result

PASS

## statement

At 20ae0b603 ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestRunAgentsDocLessons|TestAgentDocsConsistent': PASS. agents/README.md lists run-agents.md. Checked on the VM: .git/rhiza/ship/<cs>.json holds a pushed key, the built binary has rhiza ship rebase, and 9 *-landed.txtar testscripts exist as the doc says.

## Links

### Out

- verifies [run-agents-doc](../capability/rhiza/integrity/own-agent-docs/run-agents-doc.md)

## History

- `01M76NXJ53FB2GVBQPT4F4DW3R` agents/developing.md, agents/run-agents.md and agents/vm.md: the lessons for developing Rhiza, for run agents working on Rhiza, and for operating the VM
