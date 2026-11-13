# message-file-in-session-dir

decision · `01M7AAAVHCY0H138EBW8W8HGMV`

## statement

A session message is written under the folder the session runs in (its record's Dir, a verifier's or review check's checkout, when set; its Worktree otherwise), and the pointer stays relative to that folder.

## wording

the planner's choice between the two fixes the bug request names (write in the session's own working directory, or an absolute pointer), made in an autonomous plan run: an agent reading outside its working directory may be refused or prompted

## Links

### Out

- derived-from [plan-fix-rhiza-session-message](../source/plan-fix-rhiza-session-message.md)

### In

- decided-by [message-file-in-session-checkout](../capability/rhiza/orchestration/step-messages/message-file-in-session-checkout.md)

## History

- `01M7AA6QGX335J8XPJX3D6T4DJ` Fix rhiza session message writes its file into the run's worktree, but verifiers and review checks run in detached checkouts under .git/rhiza/checkouts/, so the session is told to read a relative path that does not exist where it runs. On 2026-11-13 at 10:05Z, rhiza session message 4594e9fa (verifier of run 01M79TRE3W, cwd .git/rhiza/checkouts/01M79TRE3WM2X83PZ6TBWRQC3W-daemon-holds-steps-on-stale-guide-2) wrote /home/rhiza/code/rhiza.fix-runs-started-agent/.rhiza/local/messages/4594e9fa.../0001.md; the session answered 'I couldn't read the message. .rhiza/local/messages/.../0001.md doesn't exist' and went idle again. Worked around at 10:17Z by creating the file at the same relative path inside the checkout before sending the next message. Wanted: write the message file in the session's own working directory (its checkout), or give the pointer as an absolute path.
