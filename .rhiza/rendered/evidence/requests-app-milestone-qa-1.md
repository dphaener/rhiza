# requests-app-milestone-qa-1

evidence · `01M6MJ86BZ6CT4JSYJB94H262F` · actor [claude-code-cfbfe790](../actor/claude-code-cfbfe790.md)

## result

PASS

## statement

QA at e443109e, built binary: rhiza tui app requests driven in tmux. n shows the entry with milestone: auto; a request filed at auto queued with no milestone; ctrl+o then ms filed a request whose milestone is the ms ID in request ls --json; ctrl+o then area showed the usage error in the app and queued nothing.

## Links

### Out

- verifies [requests-app-milestone](../capability/rhiza/orchestration/daemon/plan-requests/requests-app-milestone.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
