# daemon-picks-verified-3

evidence · `01M62C49F1KDPZY16W4E5EDRYV` · actor [claude-code-33dafe5d](../actor/claude-code-33dafe5d.md)

## result

PASS

## statement

At 4880e32, go test -count=1 ./internal/api -run 'TestDaemonPicks|TestDaemonHeld' passes: over a fixture with first/second milestone change sets with leaves, a leafless milestone one, a backlog one with leaves, a draft, and running/interrupted/failed runs on three more, daemon.Pick picks only first under max 2, none under max 1, and with the failed run abandoned picks failed again (max 4: first, second, failed); a pending start counts as running. go test ./internal/daemon passes (TestResumable now also skips a locked run and one with no process record younger than a minute). go test ./cmd/rhiza -run TestScripts/daemon passes daemon-picks, daemon-status, daemon-resume, daemon-lock; daemon-picks runs rhiza daemon --once with a fake Rex and sees exactly one 'rex window new ... run start <id> --permission-mode manual' while the leafless and backlog change sets stay unstarted. Read pick.go: live (running or interrupted) and failed-not-abandoned latest runs are skipped, only running runs and pending starts take slots. make fmt-check vet test pass; make self-check fails only on stale .rhiza/rendered files from earlier Evidence writes, not on code.

## Links

### Out

- verifies [daemon-picks](../capability/rhiza/orchestration/daemon/daemon-picks.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
