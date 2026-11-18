# request-notes-plan-brief-qa-1

evidence · `01M7QWM2MWC4M6PA7655RKP2DW` · actor [claude-code-059bdab0](../actor/claude-code-059bdab0.md)

## result

PASS

## statement

QA with a scratch build and fake rex/claude: queued a request with two --note values (one multi-line Markdown) plus a third via request note, rhiza daemon --once drew its run, run start --foreground --request ran the plan workflow; the planner's brief.md had ## Request notes right after ## Change set, the lead sentence and Note 1..3 verbatim in order, then ## Leaves. A plan run with no request had no Request notes section. go test ./internal/run ./internal/daemon -run Brief|Request, ./internal/api -run Request and the request-notes-brief, describe-golden, plan-requests, brief-files scripts pass.

## Links

### Out

- verifies [request-notes-plan-brief](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/request-notes-plan-brief.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
