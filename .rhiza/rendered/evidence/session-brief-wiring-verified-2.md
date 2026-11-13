# session-brief-wiring-verified-2

evidence · `01M7B1G2MGPY5ZPMA8D86FVMB5` · actor [claude-code-37074b57](../actor/claude-code-37074b57.md)

## result

PASS

## statement

Verified attempt 2 at 9ba42f285 (review-cycle-2 fix f1: brief.md lists the change set's own focus Constraints, then the others, like Decisions). Acceptance: go test -trimpath ./internal/api -run TestSessionBrief exit 0 (TestSessionBrief, TestSessionBriefUnwritable, TestSessionBriefDecisions, TestSessionBriefOwnDecisionsConstraints PASS); go test -trimpath ./cmd/rhiza -run TestScripts/brief-files exit 0. Checks block: before-fix run (HEAD~1 run_brief.go in a scratch worktree) of TestSessionBriefOwnDecisionsConstraints exit 1, failing on missing '## Constraints - No network.' as claimed; -run 'TestSessionBrief|TestArbiterBrief|TestAdvisorBrief' exit 0; go vet ./internal/api exit 0; rhiza ship generate exit 0 (0 changed); rhiza render exit 0 (no tracked changes). The full go test -timeout 8m -skip TestReadBudget ./internal/api/ hit the 8m timeout while TestWritesThroughIndexMatchFold was running, at host load average ~10, with no test failure reported. The changed code (focusStatements, briefInput) is reached only via fillSession, and the session, brief, leaves and checks tests that use it pass (-run 'Brief|Session|StartAgent|Leaves|Checks|Handback|Slice' exit 0), as does ./internal/run. The diff touches only internal/api/run_brief.go, its test and .rhiza records. No TODOs or stubs; fakes only.

## Links

### Out

- verifies [session-brief-wiring](../capability/rhiza/workflows/workflow-runs/session-brief-wiring.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
