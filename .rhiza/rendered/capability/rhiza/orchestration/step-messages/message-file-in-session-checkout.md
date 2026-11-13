# message-file-in-session-checkout

capability · `01M7AAAVJE3QMNAYAZKHG3ASHW` · parent [step-messages](../step-messages.md)

## acceptance

New tests in internal/api/session_message_test.go and internal/api/session_ask_test.go set a session record's Dir to a separate checkout folder and show: session message (and its dry run's File) puts 0001.md under <Dir>/.rhiza/local/messages/<session>/ and nothing under the Worktree's folder, with the pointer naming .rhiza/local/messages/<session>/0001.md; Deliver writes the answer there too; a session with no Dir still writes under its Worktree. go test -trimpath ./internal/api -run 'TestSessionMessage|TestDeliver|TestTypePointer' and go test -trimpath ./cmd/rhiza -run 'TestScripts/(session-message|describe-golden)' pass, describe-golden.txtar regenerated if the File description changes.

## statement

A message to a session that runs in a detached checkout (a verifier or review check, its record's Dir set) is written under that checkout's .rhiza/local/messages/<session>/, so the relative pointer names a file that exists where the session runs; this holds for rhiza session message, its dry run, and the arbiter's answers that Deliver sends, through messageFolder in internal/api/session_message.go, whose File response description and doc comments stop saying the run's worktree.

## Links

### Out

- decided-by [message-file-in-session-dir](../../../../decision/message-file-in-session-dir.md)

### In

- verifies [message-file-in-session-checkout-verified](../../../../evidence/message-file-in-session-checkout-verified.md)
- verifies [message-file-in-session-checkout-qa-1](../../../../evidence/message-file-in-session-checkout-qa-1.md)
- verifies [pr-545-fix-rhiza-session-message](../../../../evidence/pr-545-fix-rhiza-session-message.md)

## History

- `01M7AA6QGX335J8XPJX3D6T4DJ` Fix rhiza session message writes its file into the run's worktree, but verifiers and review checks run in detached checkouts under .git/rhiza/checkouts/, so the session is told to read a relative path that does not exist where it runs. On 2026-11-13 at 10:05Z, rhiza session message 4594e9fa (verifier of run 01M79TRE3W, cwd .git/rhiza/checkouts/01M79TRE3WM2X83PZ6TBWRQC3W-daemon-holds-steps-on-stale-guide-2) wrote /home/rhiza/code/rhiza.fix-runs-started-agent/.rhiza/local/messages/4594e9fa.../0001.md; the session answered 'I couldn't read the message. .rhiza/local/messages/.../0001.md doesn't exist' and went idle again. Worked around at 10:17Z by creating the file at the same relative path inside the checkout before sending the next message. Wanted: write the message file in the session's own working directory (its checkout), or give the pointer as an absolute path.
