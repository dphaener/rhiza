# session-brief-wiring

capability · `01M7AHST8TM851DP82HM2RZ85M` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/api -run TestSessionBrief shows an agent step's, a builder's and a checks step's verifier's session getting their own directory with brief.md and slice.json, slice.json equal to rhiza slice --json for the change set, the three variables set only when their file exists, and a session started outside a run getting none; go test -trimpath ./cmd/rhiza -run TestScripts/brief-files passes with a fake agent printing its RHIZA_BRIEF file, and its failure case: an unwritable runs directory fails the session start with a message naming the path and starts no agent. Touches: internal/api/run_brief.go (new), internal/api/run_brief_test.go (new), internal/api/run.go (StartAgent), internal/api/run_leaves.go, internal/api/run_checks.go, internal/api/session_run.go (launchEnv), cmd/rhiza/testdata/script/brief-files.txtar (new). Must not change prompt rendering or the workflow file format.

## statement

Before every agent session it starts for a run from a workflow step's table (an agent step, a leaves step's builder and verifier, a checks step's check, critic, fixer and verifier), the daemon fills the session's directory with brief.md built from the model and slice.json equal to rhiza slice --json for the change set at that moment, and sets RHIZA_BRIEF and RHIZA_SLICE (and RHIZA_HANDBACK when handback.md exists there) in the session's environment; a session outside a run gets none of them.

## Links

### Out

- depends-on [run-brief-files](run-brief-files.md)
- decided-by [brief-files-and-invoke-placement](../../../../decision/brief-files-and-invoke-placement.md)
- decided-by [conductor-milestone-direction](../../../../decision/conductor-milestone-direction.md)

### In

- depends-on [agent-invoke-template](../workflow-definitions/agent-invoke-template.md)
- depends-on [arbiter-advisor-brief](arbiter-advisor-brief.md)
- verifies [session-brief-wiring-verified](../../../../evidence/session-brief-wiring-verified.md)
- verifies [session-brief-wiring-qa-1](../../../../evidence/session-brief-wiring-qa-1.md)
- verifies [session-brief-wiring-verified-2](../../../../evidence/session-brief-wiring-verified-2.md)
- verifies [session-brief-wiring-qa-2](../../../../evidence/session-brief-wiring-qa-2.md)
- verifies [session-brief-wiring-verified-5](../../../../evidence/session-brief-wiring-verified-5.md)
- verifies [pr-565-every-agent-session-a-run](../../../../evidence/pr-565-every-agent-session-a-run.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
