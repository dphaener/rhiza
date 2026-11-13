# plan-fix-rhiza-session-message

source · `01M7AAACQS5Z4ERKP5BDFWHF1A`

## quotes

```json
[]
```

## raw_sha256

5b327cadf591e7fbf5a5d3bb0c306f4b63b0c12423239c4617a4bb6a4fe8ae14

## session

5fa83e6d-6fde-4e55-98fe-01a1d4a3ae5e

## summary

Autonomous plan run for a bug request; no owner interview. The intent reports that rhiza session message writes its file under the run's worktree while verifier and review-check sessions run in detached checkouts under the git common dir's rhiza/checkouts, so the relative pointer names a file that does not exist where the session runs. Code reading: messageFolder in internal/api/session_message.go joins session.Record.Worktree; the record's Dir field already holds the checkout a verifier or check session runs in (internal/session/session.go). messageFolder serves both session.message and the arbiter's Deliver (internal/api/run.go), so one fix covers both. The intent allows either writing the file in the session's own working directory or an absolute pointer; the planner chose the working directory, keeping the relative pointer, because an agent reading a path outside its working directory may be refused or prompted. The checkout's .rhiza/.gitignore ignores local/, so the file stays untracked there. Not severe: a workaround existed and no data was lost. Not breaking: no command, flag or JSON shape changes. Nothing is left open.

## title

Plan: fix rhiza session message in checkouts (autonomous)

## Links

### In

- derived-from [message-file-in-session-dir](../decision/message-file-in-session-dir.md)

## History

- `01M7AA6QGX335J8XPJX3D6T4DJ` Fix rhiza session message writes its file into the run's worktree, but verifiers and review checks run in detached checkouts under .git/rhiza/checkouts/, so the session is told to read a relative path that does not exist where it runs. On 2026-11-13 at 10:05Z, rhiza session message 4594e9fa (verifier of run 01M79TRE3W, cwd .git/rhiza/checkouts/01M79TRE3WM2X83PZ6TBWRQC3W-daemon-holds-steps-on-stale-guide-2) wrote /home/rhiza/code/rhiza.fix-runs-started-agent/.rhiza/local/messages/4594e9fa.../0001.md; the session answered 'I couldn't read the message. .rhiza/local/messages/.../0001.md doesn't exist' and went idle again. Worked around at 10:17Z by creating the file at the same relative path inside the checkout before sending the next message. Wanted: write the message file in the session's own working directory (its checkout), or give the pointer as an absolute path.
