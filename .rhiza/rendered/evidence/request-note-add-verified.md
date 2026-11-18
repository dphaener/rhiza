# request-note-add-verified

evidence · `01M7QVSZ1F8HXQKFKHHYYNDNQ3` · actor [claude-code-01ebe82a](../actor/claude-code-01ebe82a.md)

## result

PASS

## statement

Verifier in a detached checkout at e91bc0f3e re-ran the Checks block: go test -trimpath -p 2 ./internal/api/ -run Request (exit 0; -v shows TestRequestNoteAdd PASS), go test ./cmd/rhiza -run 'TestScripts/describe-golden$' with -update then without (both exit 0, golden unchanged), rhiza ship generate and rhiza render (exit 0, tree clean). TestRequestNoteAdd queues with one --note, adds a second via request.note, request ls shows both in order, a 6-char prefix resolves, --dry-run writes nothing, and drawn (Run set), unknown, ambiguous, short-prefix, empty and non-UTF-8 notes are usage errors with requests.json byte-identical. Diff writes only through daemon.UpdateRequests (dry run reads only), refuses Run!="" per request-note-while-queued, leaves rhiza request's flags unchanged, no TODOs. Smoke test of a scratch build: note by prefix, --dry-run, abc and blank note refused with exit 2.

## Links

### Out

- verifies [request-note-add](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/request-note-add.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
