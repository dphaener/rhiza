# request-note-shape

decision · `01M7QV3MX64Z7YN4PS45QXGE9D`

## reasoning

The orchestrator's real planner notes (~/orchestrator/planner-notes/*.md, 548 to 1997 bytes, numbered Markdown lists) are what a request should carry, so a note keeps its lines; a byte cap well above them keeps a pasted log out of every planner's brief

## statement

A request note is free text, possibly several lines (Markdown, as the orchestrator's planner-notes files are), stored on the request in requests.json (daemon.Request.Notes, oldest first) as given apart from trimming leading and trailing blank space and turning CRLF into LF; an empty note is refused, and a note over 8000 bytes or not valid UTF-8 is refused; rhiza request ls shows each request's notes (notes in --json, the count in text)

## wording

(assumed, would have escalated) the conductor's ruling, not asked of the owner

## Links

### Out

- derived-from [m40-request-notes-conduct-plan](../source/m40-request-notes-conduct-plan.md)

### In

- decided-by [request-note-flag](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/request-note-flag.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
