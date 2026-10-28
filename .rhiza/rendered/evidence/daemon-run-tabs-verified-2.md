# daemon-run-tabs-verified-2

evidence · `01M62AZF8NYGY2K7APW2SX97KM` · actor [claude-code-795d7b63](../actor/claude-code-795d7b63.md)

## result

PASS

## statement

At 161f5f5 (the fix adding daemon --permission-mode), go test ./cmd/rhiza -run 'TestScripts/(daemon-.*|describe-golden)$' -count=1 passes, and go test ./... passes. In daemon-picks.txtar's kept work dir, the fake Rex's rex.log has exactly one new tab per daemon --once poll: 'rex window new 01M5PWGQ8 -c <repo> --json -- <rhiza> run start 01M5PWGQ8GB4618G4RH2TV0ZAH --permission-mode manual'. That is the Built change set's ID with no --autonomous, and the leafless and backlog change sets got no tab. Under daemon --once --autonomous the line ends '--permission-mode manual --autonomous'. The daemon now also needs --permission-mode, which the QA blocker asked for: without it, or with bypassPermissions, it exits non-zero and opens no tab (asserted by '! exists rex.log'). The script's ID assertion is the regex [0-9A-Z]*, which would also match an empty ID, so I read the log by hand to confirm the ID. make self-check fails only on stale .rhiza/rendered files, not on code.

## Links

### Out

- verifies [daemon-run-tabs](../capability/rhiza/orchestration/daemon/daemon-run-tabs.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
