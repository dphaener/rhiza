# handback-clears-live-done-mark-qa-2

evidence · `01M7MG3HBK0XQ99SVZAE69JWAS` · actor [claude-code-cd0732bb](../actor/claude-code-cd0732bb.md)

## result

PASS

## statement

QA at 19080e6f6: built rhiza to a scratch path and ran a scratch workflow (agent step build with until_when=done and a handback template, command step check failing once with on.fail=handback:build) under a fake rex and a fake live claude that reports Stop, runs rhiza session done, then blocks until rex send types the hand-back. Fixed binary: the run succeeded via 'rejoining session', the raw session record at the start of visit 2 had no done mark, and the final mark came from the second, fresh session done. Control with the merge-base binary (abd7d1a89): visit 2's record still held visit 1's done mark and the final mark was that stale one. Failure cases: rex send failing ('rex could not type', run failed) and a step with no handback template ('step build sets no handback template', no message written) both left visit 1's done mark in place. Scoped tests: go test ./internal/api -run 'TestHandbackDeliver|TestSessionDone|TestSessionMessage' ok; go test ./cmd/rhiza -run 'TestScripts/(run-handback|ship-change-one-ci)$' ok.

## Links

### Out

- verifies [handback-clears-live-done-mark](../capability/rhiza/workflows/workflow-runs/run-handback/handback-clears-live-done-mark.md)

## History

- `01M7MD2K2BE981QPPF5HE8M551` Fix A rejoined conductor session keeps its first visit's done mark, so after a gate hand-back the session reads as done before it has worked the hand-back (anything closing done sessions closes it mid-visit)
