# request-notes-in-plan-brief

decision · `01M7QV3N1YAQBJMZ732P5DPQ0S`

## reasoning

brief.md is the file every step session already reads, and the plan run's ID ties it to its request with no new journal field

## statement

Every agent step session of a request's plan run (the run whose ID is the request's Run) gets a Request notes section in brief.md, right after the Change set section: one sentence saying the notes came with the request and are ideas to weigh against the code, then each note, oldest first, under its own '### Note <n>' heading, verbatim; a run that is no request's, or a request with no notes, renders no section. The notes are read from requests.json when the session's brief is written. Spawned children's briefs are their caller's and are not changed

## wording

(assumed, would have escalated) the conductor's ruling, not asked of the owner

## Links

### Out

- derived-from [m40-request-notes-conduct-plan](../source/m40-request-notes-conduct-plan.md)

### In

- decided-by [request-notes-plan-brief](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/request-notes-plan-brief.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
