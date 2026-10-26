# blocked-run-notice-tested

evidence · `01M5WGZHXYGN6K3NQR072DN0EV` · actor [claude-code-01m5szj4](../actor/claude-code-01m5szj4.md)

## command

go test ./internal/run -run TestAgent && go test ./internal/api -run 'TestRunBlocked|TestRunInfoHumanBlocked|TestRexReportEncoder'

## result

pass

## statement

Go tests with a fake Rex writer count one OSC 7501 blocked report per block, none while the run stays blocked, and a working report when the block ends.

## tests

internal/run: TestAgentBlockedAndBack, TestAgentSecondBlockNotifiesAgain, TestAgentBlockWithoutNotify; internal/api: TestRunBlockedStatusAndRexReports, TestRunBlockedWithoutRexWritesNothing, TestRexReportEncoder

## Links

### Out

- verifies [blocked-run-notice](../capability/rhiza/orchestration/notifications/blocked-run-notice.md)

## History

- `01M5SZJ46HYE2DDJ0STCDZE1WP` Blocked runs tell the owner: a run waiting on the owner sends a Rex notification, and genesis run status names the tab and the question.
