# requests-app-milestone-verified

evidence · `01M6MFYQ6ZG70B1MFHNHGE2SY3` · actor [claude-code-317bd79b](../actor/claude-code-317bd79b.md)

## result

PASS

## statement

At 12d4691c ran go test -trimpath ./internal/tui/requestsapp/... -run Milestone -v: TestMilestoneFieldDefaultsToAutoAndIsFiledWhenChosen PASS. It files a request with the milestone left at auto and gets argv [--kind feature -- dark mode] (no --milestone), then chooses milestone-25 via ctrl+o and gets [--kind feature --milestone milestone-25 -- paging]; the entry view shows 'milestone: auto' and then 'milestone: milestone-25'. Also ran the requestsapp package and cmd/rhiza unit tests (TestScripts skipped): ok, incl. TestRequestsAppArgsKeepDashLeadingText. Read the diff: RequestArgs adds --milestone only when non-empty, before '--'; api.RequestRequest has flag:"milestone"; gofmt clean.

## Links

### Out

- verifies [requests-app-milestone](../capability/rhiza/orchestration/daemon/plan-requests/requests-app-milestone.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
