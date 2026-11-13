# arbiter-advisor-brief-verified

evidence · `01M7AVQ7S53TA0E8AN9CB32N9N` · actor [claude-code-44762391](../actor/claude-code-44762391.md)

## result

PASS

## statement

At 56d9f9a9e re-ran the Checks block: go test -trimpath ./internal/api -run 'TestArbiterBrief|TestAdvisorBrief' (exit 0; TestArbiterBrief, TestAdvisorBrief, TestArbiterBriefUnwritable PASS), go test -trimpath ./internal/run/... (exit 0), rhiza ship generate (exit 0, 0 files changed), rhiza render (exit 0); also go vet on internal/api and internal/run (exit 0) and go test -trimpath ./internal/api -run 'Arbiter|Advisor|Brief|SessionStart|Session' (exit 0). Read the diff: StartArbiter/StartAdvisor call prepareKeyed, which has runSessionStart fill <commondir>/rhiza/runs/<run>/arbiter|advisor/<session id>/ (run.KeyedSessionDir) with brief.md and slice.json after the session ID is allocated and before the record is written or Rex is called. This is the layout Decision arbiter-advisor-dirs-per-session sets, answering review question r3-f2 and replacing the leaf's original <arbiter step>/arbiter-<question id> and advisor/advisor-<asker> names. The tests check RHIZA_BRIEF and RHIZA_SLICE, that slice.json equals the slice op's canonical output, that the arbiter and advisor prompts are unchanged, and that a second advisor gets its own directory and leaves the first one's files alone. The failure case (a runs path that is a file) returns an error naming the path, with no Rex call. arbiterWrites is untouched. No change outside the leaf's files beyond the session.go hook and the actor record, and no TODOs.

## Links

### Out

- verifies [arbiter-advisor-brief](../capability/rhiza/workflows/workflow-runs/arbiter-advisor-brief.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
