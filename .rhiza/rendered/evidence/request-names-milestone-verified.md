# request-names-milestone-verified

evidence · `01M6MG5Q7HD9T2KEMDZ23AYBVQ` · actor [claude-code-8b406559](../actor/claude-code-8b406559.md)

## result

PASS

## statement

At 3acfa0b0 ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(plan-requests|describe-golden)$': ok. plan-requests queues 'Add paging' --milestone ms (an intent with milestone=1) and request ls --json shows a 26-char milestone ID; --milestone area (intent with no milestone field) and --milestone nonesuch exit with usage errors ('not a live milestone intent', 'no such node') and request ls then shows No plan requests. By hand in a scratch repo with a hermetic HOME and a scratch build: the --milestone area request exited 2 and queued nothing, and request ls --json's milestone equalled rhiza show ms's node ID exactly (01M6MF5HYZVK4ZV6831VKC3PP1). go test ./internal/daemon ok; go test ./internal/api -skip TestReadBudget ok. TestReadBudget failed only on roadmap/ready what-if medians over 50ms at load average 14; the diff does not touch those read paths. The diff checks Retired, Kind==intent and the milestone field, and stores the ID in daemon.Request.Milestone and RequestInfo; describe-golden was regenerated. Note: the plain-text request ls does not print the milestone, only --json does, which the acceptance asks for.

## Links

### Out

- verifies [request-names-milestone](../capability/rhiza/orchestration/daemon/plan-requests/request-names-milestone.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
