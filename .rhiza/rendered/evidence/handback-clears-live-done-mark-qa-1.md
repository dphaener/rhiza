# handback-clears-live-done-mark-qa-1

evidence · `01M7MEQA6RDGXDAX061G26VKTY` · actor [claude-code-6ab73db0](../actor/claude-code-6ab73db0.md)

## result

PASS

## statement

QA with a scratch build (go build -trimpath -o $TMPDIR/rhiza-qa ./cmd/rhiza), fakes only (rex that backgrounds the tab command and pipes send to a fifo, a live claude that hooks Stop, says done and blocks on the fifo, a check failing once): rhiza run start --foreground on a workflow whose agent step waits for done (until_when=done) with a handback template. At hand-back the raw session record had no done mark, the message 0001.md pointer was typed (rex send + send-key enter), the second session done said 'is done' (not already), build visit 2 passed only after it, run succeeded. The same scenario on the merge-base build kept the done mark at hand-back. Two consecutive hand-backs (check failing twice) cleared the mark each time. Failure: the same step without a handback template failed the run with 'step build sets no handback template...', the raw record kept its done mark, no message file and no rex send. go test -trimpath -p 2 ./internal/api -run TestHandbackDeliver passed.

## Links

### Out

- verifies [handback-clears-live-done-mark](../capability/rhiza/workflows/workflow-runs/run-handback/handback-clears-live-done-mark.md)

## History

- `01M7MD2K2BE981QPPF5HE8M551` Fix A rejoined conductor session keeps its first visit's done mark, so after a gate hand-back the session reads as done before it has worked the hand-back (anything closing done sessions closes it mid-visit)
