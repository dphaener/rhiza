# handback-clears-live-done-mark-verified

evidence · `01M7ME6Z4AJJMFJ65NVG064DB0` · actor [claude-code-adf2b971](../actor/claude-code-adf2b971.md)

## result

PASS

## statement

At commit 05e811bf3: go test -trimpath -p 2 ./internal/api -run 'TestHandbackDeliver' -v passes (9 subtests, including the new 'a live session that said done on an earlier visit has its done mark cleared' and 'a failed hand-back leaves a live session's done mark in place'); go test -trimpath -p 2 ./internal/api/... ok; gofmt -l internal/api empty. The happy subtest covers the acceptance as written: the done mark is cleared, the hook stays idle, message 0001.md is written ending with the handback text, rex send is called (typed), AgentStatus reports SaidDone false, then session.done gives a fresh mark (not Already) and AgentStatus reports SaidDone true. The failure subtest uses a step with no handback template: the error names step build, the done mark is kept, no message is written and there are no rex calls. Diff read against must_not: only Status.Done is cleared in messageSession (asked, hook, rex, started and prompted are untouched). messageSession's only caller is Handback's live branch; resumeSession, session done/message and internal/run are unchanged. The mark is cleared only after the messageable, sender and folder checks, so a template failure, which returns earlier in Handback, leaves the mark.

## Links

### Out

- verifies [handback-clears-live-done-mark](../capability/rhiza/workflows/workflow-runs/run-handback/handback-clears-live-done-mark.md)

## History

- `01M7MD2K2BE981QPPF5HE8M551` Fix A rejoined conductor session keeps its first visit's done mark, so after a gate hand-back the session reads as done before it has worked the hand-back (anything closing done sessions closes it mid-visit)
