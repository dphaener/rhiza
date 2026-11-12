# developing-doc-qa-1

evidence · `01M7876VCYWP7ZGKXAEK2NBJ39` · actor [claude-code-fcd1aaeb](../actor/claude-code-fcd1aaeb.md)

## result

PASS

## statement

At 20ae0b603 built rhiza into the session scratchpad and ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestDevelopingDocLessons|TestAgentDocsConsistent' with GOTMPDIR/TMPDIR in the scratchpad: all PASS. agents/README.md lists developing.md. Checked the doc's lessons against the system: testscript TestScripts/apply (rhiza init) passes with the run's own GOTMPDIR, since #488 already puts run GOTMPDIR under ~/.cache/rhiza/run-tmp, so the GOTMPDIR section's reason is stale (reported as a minor finding).

## Links

### Out

- verifies [developing-doc](../capability/rhiza/integrity/own-agent-docs/developing-doc.md)

## History

- `01M76NXJ53FB2GVBQPT4F4DW3R` agents/developing.md, agents/run-agents.md and agents/vm.md: the lessons for developing Rhiza, for run agents working on Rhiza, and for operating the VM
