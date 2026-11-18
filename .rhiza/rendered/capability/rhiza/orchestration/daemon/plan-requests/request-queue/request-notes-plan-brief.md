# request-notes-plan-brief

capability · `01M7QV3NDW5YGBD04964SJCJ0G` · parent [request-queue](../request-queue.md)

## acceptance

internal/run brief tests render a BriefInput with two request notes, one multi-line (section after Change set, both notes verbatim in order under their headings) and one with none (no section, bytes unchanged from today). A new cmd/rhiza testscript request-notes-brief.txtar queues a request with two --note values, lets the daemon draw its run (fake rex, as plan-requests.txtar does), runs that plan run by hand with a fake claude, and finds both notes in the planner session's brief.md; (failure case) a plan run started for a request with no notes has no Request notes section in its brief. Run: go test -trimpath -p 2 ./internal/run/ -run Brief, go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/(request-notes-brief|brief-files)$'

## files

internal/run/brief.go, internal/run/brief_test.go, internal/api/run_brief.go, cmd/rhiza/testdata/script/request-notes-brief.txtar

## must_not

Change brief.md for a run that is no request's; change spawned children's briefs (RenderSpawnBrief); add a field to the run journal; start a process in the testscript unless it is listed on processScripts with its reason

## statement

The brief.md of every agent step session of a request's plan run shows the request's notes in a Request notes section after the Change set section, as request-notes-in-plan-brief says (each note verbatim under '### Note <n>', oldest first; run.BriefInput gains the notes; stepRunner.fillSession reads them from requests.json by the run's ID); a run that is no request's, or one whose request has no notes, renders no such section

## Links

### Out

- depends-on [request-note-flag](request-note-flag.md)
- decided-by [request-notes-in-plan-brief](../../../../../../decision/request-notes-in-plan-brief.md)

### In

- verifies [request-notes-plan-brief-verified](../../../../../../evidence/request-notes-plan-brief-verified.md)
- verifies [request-notes-plan-brief-qa-1](../../../../../../evidence/request-notes-plan-brief-qa-1.md)

## History

- `01M7NVMQER3WR98SDXHQCBXKSY` A request carries notes (rhiza request --note, repeatable, and notes added to a queued request) that the planner of the plan run it starts is shown in its brief, so ideas reach planners with no orchestrator messaging them
