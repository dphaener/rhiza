# request-note-while-queued

decision · `01M7QV3MZBRKBZ46M9ZYW1ZD43`

## reasoning

The intent says notes added to a queued request; once the run is drawn its planner may already have read its brief, so a later note would silently never reach it

## statement

rhiza request note <request> <text> adds a note to a request only while it is queued (no plan run drawn for it); a request whose plan run has started, or an unknown one, is a usage error and nothing is written

## wording

(assumed, would have escalated) the conductor's ruling, not asked of the owner

## Links

### Out

- derived-from [m40-request-notes-conduct-plan](../source/m40-request-notes-conduct-plan.md)

### In

- decided-by [request-note-add](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/request-note-add.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
