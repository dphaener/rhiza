# run-agents-doc-verified

evidence · `01M786XTD016EP86ATF32QK6K9` · actor [claude-code-48e5ee54](../actor/claude-code-48e5ee54.md)

## result

PASS

## statement

At 43b42aef3 ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestRunAgentsDocLessons|TestAgentDocsConsistent': TestRunAgentsDocLessons, TestAgentDocsConsistent and TestAgentDocsConsistentFixtures all PASS. Read the diff: agents/run-agents.md is listed in agents/README.md and covers trailers in one paragraph with Co-Authored-By, a run change set that already landed, an empty commit for findings HEAD already fixes, merging origin/main rather than rebase and force-push, and updating pushed in .git/rhiza/ship/<cs>.json after a hand push. The test checks for each of these lessons. The commit's trailers parse as one block (git interpret-trailers --parse).

## Links

### Out

- verifies [run-agents-doc](../capability/rhiza/integrity/own-agent-docs/run-agents-doc.md)

## History

- `01M76NXJ53FB2GVBQPT4F4DW3R` agents/developing.md, agents/run-agents.md and agents/vm.md: the lessons for developing Rhiza, for run agents working on Rhiza, and for operating the VM
