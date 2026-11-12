# developing-doc

capability · `01M76P0B69F8152BRFFCTW5G8X` · parent [own-agent-docs](../own-agent-docs.md)

## acceptance

A test in cmd/rhiza reads agents/developing.md and finds each lesson (testscript, scratchpad, RHIZA_ROOT, GOTMPDIR, root daemon, rhiza render); it and the agent-docs consistency test pass with go test -trimpath ./cmd/rhiza -run on their names

## statement

agents/developing.md, listed in agents/README.md, holds the lessons for any agent changing Rhiza's code: testscript gotchas ($WORK quoting, -count, chmod, stdin paths, fake binaries), build scratch binaries into the session scratchpad rather than a shared $TMPDIR path, unset RHIZA_ROOT, RHIZA_CHANGE and RHIZA_RUN before manual rhiza checks in scratch repos, override a run's GOTMPDIR so testscript $WORK stays outside the real repo, re-run web tests that see the root daemon from a scratch worktree, and run rhiza render after landed model writes

## Links

### In

- verifies [developing-doc-verified](../../../../evidence/developing-doc-verified.md)
- verifies [developing-doc-qa-1](../../../../evidence/developing-doc-qa-1.md)
- verifies [pr-497-role-docs-for-rhiza-agents](../../../../evidence/pr-497-role-docs-for-rhiza-agents.md)

## History

- `01M76NXJ53FB2GVBQPT4F4DW3R` agents/developing.md, agents/run-agents.md and agents/vm.md: the lessons for developing Rhiza, for run agents working on Rhiza, and for operating the VM
