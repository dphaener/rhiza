# daemon-run-tabs-qa-3

evidence · `01M62D0EV4F0QX6YKWZS55ATKP` · actor [claude-code-73eca481](../actor/claude-code-73eca481.md)

## result

PASS

## statement

QA cycle 3 with the built binary (359e3c1) in a scratch repo, with a fake rex that logs every call and runs only the daemon's own `rhiza run start|resume` tab command (agent tabs only logged): `rhiza daemon --once` without --permission-mode exited 2 before any tab; with --permission-mode manual the tab ran '<rhiza> run start <cs id> --permission-mode manual', which passed run start, journaled and reached build/build; with --autonomous the command was 'run start 01M62CZE11ZQDW7E1N866YV325 --permission-mode manual --autonomous'.

## Links

### Out

- verifies daemon-run-tabs (retired)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
