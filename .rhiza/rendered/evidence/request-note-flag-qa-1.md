# request-note-flag-qa-1

evidence · `01M7QWM2GEBGMEX9FB7RP671SV` · actor [claude-code-059bdab0](../actor/claude-code-059bdab0.md)

## result

PASS

## statement

QA built ./cmd/rhiza to a scratch path and, in a scratch root, ran rhiza request 'first idea' --note (CRLF, multi-line Markdown with outer blank lines) --note second --json: both notes stored oldest first in requests.json and request ls --json, lines kept, CRLF to LF, outer blank trimmed; request ls text shows [2 notes]/[1 note]; a request with no notes writes no notes key. Blank and empty --note (exit 2 'is empty'), 8001 bytes (exit 2 'is 8001 bytes, at most 8000') queue nothing; 8000 bytes accepted.

## Links

### Out

- verifies [request-note-flag](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/request-note-flag.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
