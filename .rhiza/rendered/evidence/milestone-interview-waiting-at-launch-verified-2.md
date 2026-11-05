# milestone-interview-waiting-at-launch-verified-2

evidence · `01M6N4T20VWNNSX89RVTPC0V6R` · actor [claude-code-095db8db](../actor/claude-code-095db8db.md)

## result

PASS

## statement

Verifier claude-code-095db8db at c58c8ead (attempt 2): go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/plan-milestone-interview-waiting passed: inside a fake rex (REX_BLOCK set) with fake claude, the non-autonomous milestone run's stderr carries exactly one state=blocked:kind=question:app=rhiza report whose msg decodes to 'interview waits on you in tab ...: the milestone interview waits on you', ordered working, blocked, done (before the fake agent asks anything), while the feature plan run and the autonomous milestone run carry no state=blocked. go test -trimpath -count=1 ./internal/api -run TestProgramStatusMilestoneInterview passed: milestoneInterviewLaunch is true only for a non-autonomous milestone plan's interview step, and launchBlocked encodes a blocked question naming the milestone interview, cleared by the next step's status. Read both leaf commits' diffs (4805285b, c58c8ead): the report is written only through ctx.Rex (driver inside Rex), matching the narrowed statement. go test ./internal/run and ./internal/api (skipping the TestReadBudget timing budget, which failed only on what-if latency under contention) pass; go vet and gofmt clean.

## Links

### Out

- verifies [milestone-interview-waiting-at-launch](../capability/rhiza/orchestration/notifications/milestone-interview-waiting-at-launch.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
