# milestone-interview-waiting-at-launch

capability · `01M6MYSV344P11Y2K0555WP96J` · parent [notifications](../notifications.md)

## acceptance

A new testscript cmd/rhiza/testdata/script/plan-milestone-interview-waiting.txtar runs a non-autonomous milestone plan run inside a fake rex (REX_BLOCK set) with fake claude: the run driver's Rex output carries a blocked OSC 7501 report naming the milestone interview before the agent asks anything, and a feature plan run's and an autonomous milestone run's carry none at launch; a new unit test TestProgramStatusMilestoneInterview in internal/api/run_status_test.go checks the launch report's encoding (blocked, naming the milestone interview) and that a feature plan run and an autonomous run get none; go test -trimpath ./cmd/rhiza -run TestScripts/plan-milestone-interview-waiting and go test -trimpath ./internal/api -run TestProgramStatusMilestoneInterview pass

## statement

When a milestone plan run driven inside Rex (REX_BLOCK set) launches its interview session, the run driver reports its own Rex status as blocked, waiting on the owner (OSC 7501), naming the milestone interview and its tab, at launch rather than at the agent's first question, so Rex notifies the owner that the interview waits on them; a daemon-started run's driver runs outside Rex and writes no such report, its notification being inbox-milestone-interview's daemon notification; other plan runs and autonomous runs report as before

## Links

### Out

- decided-by [milestone-interview-starts-at-once](../../../../decision/milestone-interview-starts-at-once.md)

### In

- verifies [milestone-interview-waiting-at-launch-verified](../../../../evidence/milestone-interview-waiting-at-launch-verified.md)
- verifies [milestone-interview-waiting-at-launch-qa-1](../../../../evidence/milestone-interview-waiting-at-launch-qa-1.md)
- verifies [milestone-interview-waiting-at-launch-verified-2](../../../../evidence/milestone-interview-waiting-at-launch-verified-2.md)
- verifies [pr-291-a-milestone-request](../../../../evidence/pr-291-a-milestone-request.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
