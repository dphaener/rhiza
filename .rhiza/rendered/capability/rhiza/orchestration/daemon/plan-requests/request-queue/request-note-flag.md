# request-note-flag

capability · `01M7QV3N5ZZR8C2R4THW49RGGN` · parent [request-queue](../request-queue.md)

## acceptance

A test in internal/api (request_test.go) queues a request with two --note values, one of them several lines, and reads both back in order, lines kept and outer blank space trimmed, from requests.json and request ls --json; (failure case) an empty or blank note, a note over 8000 bytes and a note that is not valid UTF-8 are each refused as usage errors and queue nothing. internal/daemon tests round-trip a Request with notes through WriteRequests/ReadRequests, and one without notes writes no notes key. cmd/rhiza's describe-golden.txtar is regenerated. Run: go test -trimpath -p 2 ./internal/api/ -run Request, go test -trimpath -p 2 ./internal/daemon/..., go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/(describe-golden|plan-requests)$'

## files

internal/daemon/requests.go, internal/daemon/requests_test.go, internal/api/request.go, internal/api/request_test.go, cmd/rhiza/testdata/script/describe-golden.txtar, docs/runs-and-daemon.md, docs/reference/cli.md (regenerated)

## must_not

Change the request file's existing keys or the shape of a request without notes; change how a request's text, ref or workflow is read; touch the Request TUI app

## statement

rhiza request takes --note, repeatable, and queues the request with those notes (daemon.Request.Notes, oldest first, normalised and bounded as request-note-shape says: multi-line kept); rhiza request ls lists each request's notes (notes array in --json, absent when none; '[N notes]' in text), and docs/runs-and-daemon.md says so

## Links

### Out

- decided-by [request-note-shape](../../../../../../decision/request-note-shape.md)

### In

- depends-on [request-note-add](request-note-add.md)
- depends-on [request-notes-plan-brief](request-notes-plan-brief.md)
- verifies [request-note-flag-verified](../../../../../../evidence/request-note-flag-verified.md)
- verifies [request-note-flag-qa-1](../../../../../../evidence/request-note-flag-qa-1.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
