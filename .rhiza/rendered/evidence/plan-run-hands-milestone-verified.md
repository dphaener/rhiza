# plan-run-hands-milestone-verified

evidence · `01M6MHTSCW9H01CVDCZWEWCFQT` · actor [claude-code-649aa7d5](../actor/claude-code-649aa7d5.md)

## result

PASS

## statement

At fb839bf1, go test -trimpath ./cmd/rhiza -run TestScripts/plan-run-milestone-env passes: the daemon starts plan runs for a request queued with --milestone ms and one without; run start --request for the first launches a fake claude whose env file holds exactly RHIZA_MILESTONE=<ms intent ID>, and for the second the file is empty even though the outer env sets RHIZA_MILESTONE= (launchEnv unsets it). The diff carries the request's Milestone through requestRun into the run's started entry (run.Started/State.Milestone), and StartAgent reads it into the session record, which launchEnv sets or unsets. The request/plan testscripts, internal/api (with the updated unset lists in session_test.go), internal/run and internal/session also pass.

## Links

### Out

- verifies [plan-run-hands-milestone](../capability/rhiza/orchestration/daemon/plan-requests/plan-run-hands-milestone.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
