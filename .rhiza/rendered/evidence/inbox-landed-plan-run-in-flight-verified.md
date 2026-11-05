# inbox-landed-plan-run-in-flight-verified

evidence · `01M6MQNGPBK0Y4QW76G6QT2W9Z` · actor [claude-code-a7fb8d19](../actor/claude-code-a7fb8d19.md)

## result

PASS

## statement

At 5c3750bd (verifier, read the diff): internal/api/inbox.go's gone() now treats a landed change set as not gone for a plan-workflow run (or a non-plan run that entered land), so the skipped and failed/gave-up branches follow. New cmd/rhiza/testdata/script/inbox-landed-plan-run.txtar runs rhiza inbox --json over fixture journals: R-PLAN (plan, interrupted, landed) is absent; G-PLAN (plan, landed, interrupted after daemon-resumed) is gave-up with next rhiza run resume G-PLAN; F-PLAN (plan, landed, failed) is failed with next rhiza run resume --retry F-PLAN; N-PLAN (ship-change on landed, never entered land) and W-PLAN (plan on withdrawn) stay skipped with rhiza run abandon; exactly 4 items. Ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(inbox-landed-plan-run|inbox.*|.*unaccepted.*)$' (covers inbox-landed-resumable, the inbox-omits-resumable-landed-run test) and -run 'TestScripts/(daemon-resume.*|describe-golden)$' (the unaccepted-run-in-inbox evidence's scripts): all ok. go test ./internal/api -run Inbox: ok; the full internal/api package passed apart from TestReadBudget, a load-sensitive timing budget unrelated to inbox.go (CI runs it).

## Links

### Out

- verifies [inbox-landed-plan-run-in-flight](../capability/rhiza/orchestration/inbox/inbox-landed-plan-run-in-flight.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
