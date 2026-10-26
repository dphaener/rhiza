# run-blocked-state-tested

evidence · `01M5WGZHTDE1S52HCSR0J9DFYX` · actor [claude-code-01m5szj4](../actor/claude-code-01m5szj4.md)

## command

go test ./internal/run -run TestAgent && go test ./internal/api -run 'TestRunBlocked|TestRunInfoHumanBlocked|TestRexReportEncoder'

## result

pass

## statement

Go tests drive a fake agent session to blocked and back, and a run whose session record is blocked reads blocked in run status.

## tests

internal/run: TestAgentBlockedAndBack, TestAgentResumeWhileBlocked, TestAgentBlockedThenExit, TestAgentBlockedThenRunnerError, TestAgentSecondBlockNotifiesAgain; internal/api: TestRunBlockedStatusAndRexReports

## Links

### Out

- verifies [run-blocked-state](../capability/rhiza/orchestration/notifications/run-blocked-state.md)

## History

- `01M5SZJ46HYE2DDJ0STCDZE1WP` Blocked runs tell the owner: a run waiting on the owner sends a Rex notification, and genesis run status names the tab and the question.
