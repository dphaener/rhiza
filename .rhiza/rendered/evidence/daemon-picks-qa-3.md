# daemon-picks-qa-3

evidence · `01M62D0ECTGRQGTMDSPDJ2MYRE` · actor [claude-code-73eca481](../actor/claude-code-73eca481.md)

## result

PASS

## statement

QA cycle 3 with the built binary (359e3c1) in a scratch repo, with a fake rex that logs every call and runs only the daemon's own `rhiza run start|resume` tab command (agent tabs only logged): one accepted milestone change set with a leaf was started once (one 'run start <id>' tab); a change set whose run start failed before journaling (its worktree path blocked: exit 3) got exactly one tab over an 85s daemon, and daemon status listed it under 'Not journaled' in text and --json; with a run of A live and an interrupted one, the next accepted change set was still started (interrupted takes no slot).

## Links

### Out

- verifies [daemon-picks](../capability/rhiza/orchestration/daemon/daemon-picks.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
