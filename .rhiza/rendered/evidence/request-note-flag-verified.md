# request-note-flag-verified

evidence · `01M7QVH19WV220TCCS9A1CT73J` · actor [claude-code-311d2ba1](../actor/claude-code-311d2ba1.md)

## result

PASS

## statement

Verified commit b950d63f1 in a detached checkout. Every Checks: command exited 0 as written: go test -trimpath -p 2 ./internal/api/ -run Request (TestRequestNotes passes, also with -count=1); go test -trimpath -p 2 ./internal/daemon/... (TestRequestsNotesRoundTrip, TestNormalizeRequestNote, TestRequestsRoundTrip now uses reflect.DeepEqual); describe-golden -update leaves no diff; go test ./cmd/rhiza -run 'TestScripts/(describe-golden|plan-requests)$'; rhiza ship generate and rhiza render change nothing. The diff stays within the leaf's files plus the rendered actor and log files. RequestRequest.Notes is a repeatable --note normalised by daemon.NormalizeRequestNote (CRLF to LF, outer blank space trimmed, lines kept; empty or blank, over 8000 bytes and invalid UTF-8 refused as Usagef '--note N ...'). daemon.Request.Notes is omitempty, so a request without notes keeps its old shape. request ls has notes in --json and [N note(s)] in text, and docs/runs-and-daemon.md and cli.md cover it. A scratch-binary smoke test confirmed two notes queued in order with lines kept, and a blank or 8001-byte note refused with exit 2 and nothing queued. Pre-existing gap outside this leaf: invalid UTF-8 given on the CLI is turned into U+FFFD by the argv-to-JSON frame (cmd/rhiza/argv.go) before validation, so it is stored rather than refused, the same as for --ref; the api-level refusal the acceptance names is tested and passes.

## Links

### Out

- verifies [request-note-flag](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/request-note-flag.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
