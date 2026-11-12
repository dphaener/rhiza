# developing-doc-verified

evidence · `01M786VH6VGVCEV1NEDCSAYB22` · actor [claude-code-5953f03e](../actor/claude-code-5953f03e.md)

## result

PASS

## statement

At fbea4c40b ran `go test -trimpath -count=1 ./cmd/rhiza -run 'TestDevelopingDocLessons|TestAgentDocsConsistent' -v` (GOTMPDIR/TMPDIR overridden to the scratchpad): TestDevelopingDocLessons, TestAgentDocsConsistent and TestAgentDocsConsistentFixtures PASS. Read the diff: agents/developing.md covers testscript gotchas ($WORK quoting, -count, chmod, stdin paths, fake binaries), scratch binaries in the session scratchpad, unsetting RHIZA_ROOT/RHIZA_CHANGE/RHIZA_RUN, overriding GOTMPDIR, re-running root-daemon web tests from a scratch worktree, and rhiza render after landed writes; agents/README.md indexes it; the test checks each lesson; gofmt clean; the commit stages only its own actor record.

## Links

### Out

- verifies [developing-doc](../capability/rhiza/integrity/own-agent-docs/developing-doc.md)

## History

- `01M76NXJ53FB2GVBQPT4F4DW3R` agents/developing.md, agents/run-agents.md and agents/vm.md: the lessons for developing Rhiza, for run agents working on Rhiza, and for operating the VM
