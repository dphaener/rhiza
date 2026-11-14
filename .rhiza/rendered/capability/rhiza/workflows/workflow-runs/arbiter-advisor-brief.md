# arbiter-advisor-brief

capability · `01M7AJZ8CQ9MFGF06EWDVC1JAJ` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/api -run 'TestArbiterBrief|TestAdvisorBrief' shows an arbiter session and an advisor session each getting its own directory named as stated, holding brief.md and slice.json equal to rhiza slice --json for the run's change set, RHIZA_BRIEF and RHIZA_SLICE set to those files, and their prompts unchanged; failure case: an unwritable runs directory fails the arbiter's start with a message naming the path and starts no session. Touches: internal/api/run.go (StartArbiter, StartAdvisor), internal/api/run_brief.go, internal/api/arbiter_test.go. Must not change the arbiter or advisor prompt or what an arbiter writes.

## statement

The arbiter and advisor sessions a run starts get their context in files as every other run session does: before starting one, the daemon fills its own directory, <commondir>/rhiza/runs/<run id>/<arbiter step>/arbiter-<question id>/ for an arbiter and <commondir>/rhiza/runs/<run id>/advisor/advisor-<asking session id>/ for an advisor, with brief.md and slice.json for the run's change set and sets RHIZA_BRIEF and RHIZA_SLICE in its environment; their prompts are unchanged.

## Links

### Out

- depends-on [session-brief-wiring](session-brief-wiring.md)
- decided-by [brief-files-and-invoke-placement](../../../../decision/brief-files-and-invoke-placement.md)
- decided-by [conductor-milestone-direction](../../../../decision/conductor-milestone-direction.md)

### In

- blocks [review-01m7ah6rhkmzrahgjkjfm4vtnt-r3-f2](../../../../question/review-01m7ah6rhkmzrahgjkjfm4vtnt-r3-f2.md)
- verifies [arbiter-advisor-brief-verified](../../../../evidence/arbiter-advisor-brief-verified.md)
- verifies [arbiter-advisor-brief-qa-1](../../../../evidence/arbiter-advisor-brief-qa-1.md)
- verifies [arbiter-advisor-brief-qa-2](../../../../evidence/arbiter-advisor-brief-qa-2.md)
- verifies [pr-565-every-agent-session-a-run](../../../../evidence/pr-565-every-agent-session-a-run.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
