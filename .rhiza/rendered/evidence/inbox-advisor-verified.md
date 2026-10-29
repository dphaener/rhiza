# inbox-advisor-verified

evidence · `01M62JSJ31TADVWCM3QGWKVXA7` · actor [claude-code-2088e925](../actor/claude-code-2088e925.md)

## result

PASS

## statement

At 25ddbfb ran go test ./cmd/rhiza -run TestScripts/inbox-advisor -count=3: ok. The script uses a fake claude adapter whose agent blocks on AskUserQuestion via the real hook path. (a) It asserts the advisor was launched once with prompt 'Recommend an answer to this question the owner was asked: Which colour? -- follow rhiza guide advisor.', one advisor and one recommended journal entry (text Use blue, detail Blue is the brand colour), rhiza inbox and inbox --json showing question, recommendation and reasons, no rex send/send-key to the asker (typed.log absent), no asked/answered/arbiter entries (grep pattern checked to match on this macOS grep), no Decision in the run worktree, and run status blocked until the owner answers. (b) A fake advisor exiting 5 leaves the inbox question with no recommendation (text and JSON) and the run blocked. Read the diff: engine.block launches StartAdvisor without waiting only for non-arbitrated question blocks; launch failure is only logged; run recommend only journals, refuses non-advisor sessions and ended runs; inbox attaches AdviceFor the block. make fmt-check, vet and test pass; make self-check fails only on stale .rhiza/rendered actor files written by this build run's own sessions, not by the leaf.

## Links

### Out

- verifies [inbox-advisor](../capability/rhiza/orchestration/inbox/inbox-advisor.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
