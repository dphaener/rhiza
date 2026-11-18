# request-note-add

capability · `01M7QV3N9F6ZF8ZZZWAH3685DT` · parent [request-queue](../request-queue.md)

## acceptance

A test in internal/api queues a request with one --note, adds a second with request.note, and request ls --json shows both in order; a prefix resolves; (failure case) a request whose plan run has started (Run set), an unknown ID, an ambiguous prefix and an empty note are each usage errors that leave requests.json unchanged. describe-golden.txtar is regenerated. Run: go test -trimpath -p 2 ./internal/api/ -run Request, go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/describe-golden$'

## files

internal/api/request.go (or a new internal/api/request_note.go), internal/api/request_test.go, cmd/rhiza/testdata/script/describe-golden.txtar, docs/runs-and-daemon.md, docs/reference/cli.md (regenerated)

## must_not

Add a note to a request whose plan run has been drawn; write requests.json outside daemon.UpdateRequests; change rhiza request's own flags

## statement

rhiza request note <request> <text> (operation request.note) appends a note to a queued request under the queue's lock (daemon.UpdateRequests), the request named by its full ID or a unique prefix of 4+ characters, the note normalised and bounded as request-note-shape says; --dry-run writes nothing; it prints the request with its notes

## Links

### Out

- depends-on [request-note-flag](request-note-flag.md)
- decided-by [request-note-while-queued](../../../../../../decision/request-note-while-queued.md)

### In

- verifies [request-note-add-verified](../../../../../../evidence/request-note-add-verified.md)
- verifies [request-note-add-qa-1](../../../../../../evidence/request-note-add-qa-1.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
