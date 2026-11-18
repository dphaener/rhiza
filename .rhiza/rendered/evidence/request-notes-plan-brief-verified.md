# request-notes-plan-brief-verified

evidence · `01M7QW4X0R0N9AAQDPR3TF1463` · actor [claude-code-30803065](../actor/claude-code-30803065.md)

## result

PASS

## statement

In a detached checkout at 25c244720 (capped scopes, RHIZA_/SHIP_/REX_ unset) re-ran the Checks: block: go test -trimpath -p 2 ./internal/run/ -run Brief exit 0; go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/(request-notes-brief|brief-files)$' exit 0 (both PASS); go test -trimpath -p 2 ./internal/api -run 'Brief|Session' exit 0; go vet ./internal/api ./internal/run ./cmd/rhiza exit 0; rhiza ship generate (scratch build) 0 changed; rhiza render 0 written; gofmt clean, no TODOs. Diff read: RenderBrief writes a Request notes section (lead sentence, '### Note <n>', verbatim, in order) right after Change set and before Previous attempt; none renders nothing. TestBriefRequestNotes pins the no-notes brief bytes exactly and checks a multi-line note verbatim and placement. stepRunner.requestNotes reads requests.json by the run's ID; a missing queue is an empty queue (no warning), an unreadable one warns on stderr. request-notes-brief.txtar queues a request with two --note values, the daemon draws both runs (fake rex), runs the plan run by hand with a fake claude and cmps the planner session's interview-1/agent/brief.md, finding both notes in order between Change set and Leaves; the no-notes request's brief has no section; brief-files (no request) still passes. The processScripts entry is needed: the script runs rhiza through exec sh -c several times. RenderSpawnBrief and the run journal are untouched, per Decision request-notes-in-plan-brief.

## Links

### Out

- verifies [request-notes-plan-brief](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/request-notes-plan-brief.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
