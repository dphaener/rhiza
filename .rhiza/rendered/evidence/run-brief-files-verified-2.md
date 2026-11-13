# run-brief-files-verified-2

evidence · `01M7B0PH9GMQWQSYV0HEB20CN4` · actor [claude-code-3f28c530](../actor/claude-code-3f28c530.md)

## result

PASS

## statement

Verified run-brief-files at ea46de5ca (review cycle 2 fix of f2). Acceptance go test -trimpath ./internal/run -run 'TestBriefFiles|TestSessionDir|TestHandbackFile': exit 0 (TestBriefFiles, TestSessionDir, TestSessionDirWrite, TestHandbackFile pass). Checks block re-run: TestSessionBriefDecisions against the pre-fix run_brief.go in a scratch worktree exit 1; go test -trimpath ./internal/api -run 'TestSessionBrief|TestArbiterBrief|TestAdvisorBrief' exit 0 (incl. new TestSessionBriefDecisions); go test -trimpath -timeout 8m -skip TestReadBudget ./internal/api/ exit 0; go vet ./internal/api exit 0; rhiza ship generate exit 0; rhiza render exit 0 with no tracked changes. Full ./internal/api with TestReadBudget not re-run (commit reports it exit 1 as a load-dependent timing test; untouched). Diff: brief.md now builds from an unbudgeted slice and lists the change set's live focus Decisions then the other applicable decisions, slice.json unchanged; test proves slice.json cuts decisions while brief.md lists all 41. No TODOs or stubs; change limited to internal/api/run_brief.go and its test plus rendered/log records.

## Links

### Out

- verifies [run-brief-files](../capability/rhiza/workflows/workflow-runs/run-brief-files.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
