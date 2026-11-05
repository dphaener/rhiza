# run-finds-sessions-ended-unseen-verified-2

evidence · `01M6MBNA0NX86WQ9SKJZQMSXRZ` · actor [claude-code-f2d0eecd](../actor/claude-code-f2d0eecd.md)

## result

PASS

## statement

Verifier, at commit 6d02f6c4 (review fixer for quality f1). I read the diff: a poll (AgentStatus) now leaves a SessionEnd mark younger than hookEndGrace (1 minute) unjudged, a rejoin still judges it at once (grace 0), and session.run no longer overwrites an exit already on the record. That closes the quality finding: a live session ending normally is no longer judged lost mid-poll, and a lost exit the run acted on stays. I ran: (1) go test ./cmd/rhiza -run 'TestScripts/run-resume-lost' -count=1. Both pass. run-resume-lost-fixer.txtar resumes a review run whose change-set fixer has (a) a SessionEnd hook report and no exit, and (b) a fake rex block-not-found answer. In each the record gets state lost with the reason ('SessionEnd without an exit report' / 'Rex block not found'), review-fix-judged shows fixer attempt 1 failed with the lost detail, fixer attempt 2 passes, alpha re-runs and the run succeeds. (2) go test ./internal/api -run 'Lost|Rejoin|HookEnd|SessionRunKeeps|FreshHook' -count=1 -v: all 8 pass. They include TestPollHookEndedIsLostWithoutRex (a mark past the grace is lost on the next poll with no rex call), the new TestPollFreshHookEndNotJudged (a fresh mark is not judged; the later exit reads as a pass, no rex call) and TestSessionRunKeepsLostExit. (3) go test ./internal/run -run 'Rejoin|Lost': 7 pass (engine rejoin of agent steps, leaves, checks, arbiters; lost session; rex unreachable). (4) make fmt-check vet: clean. make test: every package passes except internal/api TestReadBudget, a read-timing budget (roadmap/ready what-if medians just over 50ms) measured at load average about 9 on 4 cores. The diff does not touch the read path. make self-check fails only on stale .rhiza/rendered actor and index files, which the run renders at close-out, not on this code.

## Links

### Out

- verifies [run-finds-sessions-ended-unseen](../capability/rhiza/workflows/workflow-runs/run-finds-sessions-ended-unseen.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
