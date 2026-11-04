# run-finds-sessions-ended-unseen-verified

evidence · `01M6M8QMZWGPTKNC3CX998BWJR` · actor [claude-code-fed4723b](../actor/claude-code-fed4723b.md)

## result

PASS

## statement

Verified at d08e4ee5. I read the diffs of 7fd9c685 and d08e4ee5. `go test ./cmd/rhiza -run TestScripts/run-resume-lost -count=1` passes. run-resume-lost-fixer.txtar resumes a review run (a checks step with a change-set fixer) whose fixer session ended unseen. It has two cases. In the hook case the record has a SessionEnd hook report and no exit. The record gets state lost with the detail 'SessionEnd without an exit report', and no rex block inspect is made. In the block case the fake rex answers 'block not found'. The record gets state lost with the detail 'Rex block not found', and exactly one block inspect is made. In both cases review-fix-judged attempt 1 is outcome fail with detail 'lost: ...', fixer attempt 2 passes with a commit, check alpha re-runs in cycle 2, and the run succeeds. internal/api TestPollHookEndedIsLostWithoutRex passes: a record that gains the SessionEnd mark mid-run is judged Lost/fail on the next AgentStatus poll, with no rex call. Rejoin/Lost tests for unreachable-Rex retry and exit-is-final also pass. The engine tests TestAgentRejoinOnResumeOnly and TestAgentLostSession pass. make fmt-check, vet and the full make test pass (a second full run hit a TestReadBudget timing miss under machine load; the diff does not touch the read path, and the test passed when re-run alone). make self-check fails only on stale .rhiza/rendered files from the builders' actor records, which close-out's rhiza render regenerates.

## Links

### Out

- verifies [run-finds-sessions-ended-unseen](../capability/rhiza/workflows/workflow-runs/run-finds-sessions-ended-unseen.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
