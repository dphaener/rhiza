# message-file-in-session-checkout-verified

evidence · `01M7ACHYBV9YNDFQFK9CNM734C` · actor [claude-code-d05d726a](../actor/claude-code-d05d726a.md)

## result

PASS

## statement

At f44a77bbb re-ran the Checks block: go test -trimpath ./internal/api -run 'TestSessionMessage|TestDeliver|TestTypePointer' (exit 0, also -count=1 with TestSessionMessageInSessionCheckout and TestDeliverInSessionCheckout passing), go test -trimpath ./cmd/rhiza -run 'TestScripts/(session-message|describe-golden)' (exit 0), rhiza ship generate (exit 0, nothing changed), rhiza render (exit 0). Diff read: messageFolder uses rec.Dir when set else rec.Worktree, used by session.message (incl. dry run) and Deliver in run.go; new tests set Dir to a separate folder and check 0001.md lands under <Dir>/.rhiza/local/messages/<session>/, not under the worktree, with the pointer naming the relative path; existing tests cover the no-Dir worktree case; File desc, summary, describe-golden and cli.md updated; no TODOs, nothing outside the leaf.

## Links

### Out

- verifies [message-file-in-session-checkout](../capability/rhiza/orchestration/step-messages/message-file-in-session-checkout.md)

## History

- `01M7AA6QGX335J8XPJX3D6T4DJ` Fix rhiza session message writes its file into the run's worktree, but verifiers and review checks run in detached checkouts under .git/rhiza/checkouts/, so the session is told to read a relative path that does not exist where it runs. On 2026-11-13 at 10:05Z, rhiza session message 4594e9fa (verifier of run 01M79TRE3W, cwd .git/rhiza/checkouts/01M79TRE3WM2X83PZ6TBWRQC3W-daemon-holds-steps-on-stale-guide-2) wrote /home/rhiza/code/rhiza.fix-runs-started-agent/.rhiza/local/messages/4594e9fa.../0001.md; the session answered 'I couldn't read the message. .rhiza/local/messages/.../0001.md doesn't exist' and went idle again. Worked around at 10:17Z by creating the file at the same relative path inside the checkout before sending the next message. Wanted: write the message file in the session's own working directory (its checkout), or give the pointer as an absolute path.
