# inbox-app-actions-verified-2

evidence · `01M675SM66BK7D875VQ2JPG02F` · actor [claude-code-aed77667](../actor/claude-code-aed77667.md)

## result

PASS

## statement

At 92739bb (review fix, attempt 2) ran go test ./internal/tui/inboxapp -count=1 -v: all 15 tests PASS, with a fake Runner and a fake Exec, no real Rex or agent. o/enter calls inbox.open with the item ID (TestOpenOnEnterAndO); a then typed text then enter calls session.message to the asking session with the text, and t sends the recommendation (TestAnswerSendsTypedText, TestTakeRecommendation); R on a failed run execs rex window new ... -- rhiza run resume <run> --retry (TestRetryOpensTab); x then y calls run.abandon, x then n calls nothing (TestAbandonNeedsY). Per the confirmed QA major (approve/request-changes froze the app by continuing the run in-process), y and c now open a Rex tab running rhiza run approve <run> and rhiza run request-changes - <run> with the typed text as the one finding on stdin, instead of calling inbox.approve/inbox.request-changes in-process; inbox.approve is only item-to-run lookup then runRunApprove (internal/api/inbox.go:427), so the same handler runs. TestApproveAndRequestChangesOpenTabs asserts the exact argv and runs the tab's sh command with echo to check the findings JSON carries 'add tests'. TestModeActsOnItsItem covers the confirmed quality major: a mode is pinned to its item, a tick that drops it cancels the mode with no call, and an answer goes to the entered-on question after the cursor moves. make fmt-check vet test exit 0; make self-check fails only on stale .rhiza/rendered/ files from the run's own actor writes, not on this commit.

## Links

### Out

- verifies inbox-app-actions (retired)

## History

- `01M6732DJSQ14B2WTZYZ9MPRA5` An inbox app lists what waits on the owner, questions, sign-offs and failed runs, and acts on each in place: answer, approve, request changes, retry, abandon or open its tab
