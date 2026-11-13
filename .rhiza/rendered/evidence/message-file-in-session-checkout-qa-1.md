# message-file-in-session-checkout-qa-1

evidence · `01M7ACS5P1GN4BRWTMJ0G6TR2G` · actor [claude-code-d32f1672](../actor/claude-code-d32f1672.md)

## result

PASS

## statement

QA at the branch head: built with go build -trimpath to a scratch path and ran it in a scratch repo with a worktree and a separate detached checkout, a hand-written session record with dir set to the checkout. rhiza session message --dry-run --json named <checkout>/.rhiza/local/messages/<session>/0001.md; two real sends wrote 0001.md and 0002.md there and nothing under the worktree or main checkout, the queued pointer named .rhiza/local/messages/<session>/0002.md and that relative path resolved from the checkout, which stayed git-clean. A record without dir still wrote under its worktree. Sending from a subdirectory and through rhiza api call session.message gave the same checkout path. Failures stayed clean: exited session exit 3, blank text exit 2, unknown session exit 4. Deliver path covered by go test -trimpath ./internal/api -run 'TestSessionMessage|TestDeliver|TestTypePointer' -count=1 (ok) and TestScripts/(session-message|describe-golden) (ok).

## Links

### Out

- verifies [message-file-in-session-checkout](../capability/rhiza/orchestration/step-messages/message-file-in-session-checkout.md)

## History

- `01M7AA6QGX335J8XPJX3D6T4DJ` Fix rhiza session message writes its file into the run's worktree, but verifiers and review checks run in detached checkouts under .git/rhiza/checkouts/, so the session is told to read a relative path that does not exist where it runs. On 2026-11-13 at 10:05Z, rhiza session message 4594e9fa (verifier of run 01M79TRE3W, cwd .git/rhiza/checkouts/01M79TRE3WM2X83PZ6TBWRQC3W-daemon-holds-steps-on-stale-guide-2) wrote /home/rhiza/code/rhiza.fix-runs-started-agent/.rhiza/local/messages/4594e9fa.../0001.md; the session answered 'I couldn't read the message. .rhiza/local/messages/.../0001.md doesn't exist' and went idle again. Worked around at 10:17Z by creating the file at the same relative path inside the checkout before sending the next message. Wanted: write the message file in the session's own working directory (its checkout), or give the pointer as an absolute path.
