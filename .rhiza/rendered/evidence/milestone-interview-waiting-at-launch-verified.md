# milestone-interview-waiting-at-launch-verified

evidence · `01M6N22HM1YWBQDYJE97EF1VR6` · actor [claude-code-4ebb2878](../actor/claude-code-4ebb2878.md)

## result

PASS

## statement

At commit 4805285b, verifier claude-code-4ebb2878 ran go test -trimpath ./cmd/rhiza -run TestScripts/plan-milestone-interview-waiting (pass): with fake rex and fake claude (which only exits), the non-autonomous milestone plan run's interview tab output carries exactly one state=blocked:kind=question:app=rhiza OSC 7501 report whose decoded msg names 'the milestone interview waits on you', ordered working, blocked, done; the feature plan run and the autonomous milestone run carry no blocked report. Ran go test -trimpath ./internal/api -run TestProgramStatusMilestoneInterview -v (PASS): it checks the launch decision (milestone yes; feature, autonomous flag, autonomous journal and other step no) and the report's encoding (blocked, question, app rhiza, msg naming the milestone interview), cleared by the next step's status. Also, as regression checks: the whole internal/api and internal/run packages, the plan/rex testscripts, gofmt and go vet all pass. Read the diff: a MilestonePlan flag on the run's Started journal entry, set from the request's --milestone; stepRunner.StartAgent calls a launched hook for the interview step only when the run is not autonomous; runRex.launchBlocked writes the blocked notice, and the next non-interview status drops it.

## Links

### Out

- verifies [milestone-interview-waiting-at-launch](../capability/rhiza/orchestration/notifications/milestone-interview-waiting-at-launch.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
