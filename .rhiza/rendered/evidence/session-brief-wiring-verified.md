# session-brief-wiring-verified

evidence · `01M7ATN448WQFSA6JZQ43ABATP` · actor [claude-code-17419b89](../actor/claude-code-17419b89.md)

## result

PASS

## statement

At 41d1705de re-ran the Checks: block: go test -trimpath ./internal/api -run TestSessionBrief (TestSessionBrief and TestSessionBriefUnwritable PASS), go test -trimpath ./cmd/rhiza -run 'TestScripts/brief-files$' (ok), every script reading .git/rhiza/runs via the grep-built regex (ok, 358s, so the attempt-1 run-nested-workflow-session failure is fixed), rhiza ship generate and rhiza render (exit 0); also TestSessionRun/TestStartCheck/TestStartFix and internal/session (ok). Diff read: prepareSession wired into StartAgent, StartLeaf (builder/verifier incl. fix cycles), StartCheck and StartFix; slice.json is canon-marshalled runSlice output, which the test compares to rhiza slice --json; launchEnv sets RHIZA_BRIEF/RHIZA_SLICE only with a Brief dir and RHIZA_HANDBACK only when handback.md exists, and unsets all three outside a run; the unwritable runs path fails the start naming the path, with no rex call and no agent (unit test and brief-files.txtar). Prompt rendering and the workflow format are unchanged; no TODOs or stubs. Extra files (session.go request field, session.Record.Brief, session_test unset lists) are needed plumbing.

## Links

### Out

- verifies [session-brief-wiring](../capability/rhiza/workflows/workflow-runs/session-brief-wiring.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
