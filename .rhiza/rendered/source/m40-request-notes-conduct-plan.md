# m40-request-notes-conduct-plan

source · `01M7QV3MVZ69652ERAK97VQD1W`

## summary

Autonomous conduct run 01M7QTZWVQ: no owner to ask. The conductor read the intent m40-request-notes (measure: a test queues a request with two notes and finds both in the plan run's brief) and the code. A plan request lives in .git/rhiza/requests.json (internal/daemon/requests.go, daemon.Request), queued by rhiza request (internal/api/request.go) and listed by rhiza request ls. The daemon draws the plan run's ID onto the request (Request.Run) before run start --request takes it, so a plan run finds its request by its own run ID. Every agent step session's brief.md is rendered by run.RenderBrief (internal/run/brief.go) from run.BriefInput, filled by stepRunner.fillSession/briefInput (internal/api/run_brief.go). Repeatable flags are []string fields with a flag tag (as change new --source). The conductor ruled the open choices as Decisions marked assumed. The orchestrator then relayed machine rules for the run and pointed to its planner-notes files (~/orchestrator/planner-notes/*.md) as real examples of the notes a request should carry: multi-line Markdown of 548 to 1997 bytes; the note shape and brief section were ruled to keep such notes whole.

## title

Conductor's plan of change set 01M7NVMQ (a request carries notes)

## Links

### In

- derived-from [request-note-shape](../decision/request-note-shape.md)
- derived-from [request-note-while-queued](../decision/request-note-while-queued.md)
- derived-from [request-notes-in-plan-brief](../decision/request-notes-in-plan-brief.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
