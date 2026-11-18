# request-note-add-qa-1

evidence · `01M7QWM2JMN5N5NEASG8X5D9EG` · actor [claude-code-059bdab0](../actor/claude-code-059bdab0.md)

## result

PASS

## statement

QA with a scratch build: rhiza request note <8-char prefix> 'third\r\nline two' appended a third note (printed with all notes, CRLF to LF); --dry-run --json showed the would-be note and requests.json stayed unchanged; unknown ID, 3-char prefix, ambiguous prefix, blank and 9000-byte notes, missing args and a drawn request (Run set, human and --json) were each usage errors (exit 2) with requests.json byte-identical; after a real daemon --once drew the run, note was refused as no longer queued.

## Links

### Out

- verifies [request-note-add](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/request-note-add.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
