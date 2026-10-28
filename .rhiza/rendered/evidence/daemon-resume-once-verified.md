# daemon-resume-once-verified

evidence · `01M62984CE1M7PV5S2GPVWVSE4` · actor [claude-code-e62e61b5](../actor/claude-code-e62e61b5.md)

## result

PASS

## statement

At f6699ba ran go test ./cmd/rhiza -run 'TestScripts/(daemon|describe-golden)' -count=1 (ok) and go test ./internal/daemon ./internal/run ./internal/api (ok). daemon-resume.txtar fixtures an unfinished journal INTERRUPTED1 whose process entry names a dead pid on this host and a FAILED00001 journal finished failed; with a fake rex logging its args, the first rhiza daemon --once prints 'Resumed run INTERRUPTED1', the rex log has exactly one 'rex window new ... run resume INTERRUPTED1' line and none naming FAILED00001, and INTERRUPTED1's journal gains one daemon-resumed entry; the second rhiza daemon --once --json (run still interrupted) reports resumed [] and the log still holds one run resume. Read the diff: daemon.Resumable picks runs that are LiveInterrupted and not DaemonResumed (unit test covers interrupted/again/failed), the poll journals daemon-resumed before opening the tab. make check passes fmt, vet and tests; self-check fails only on rendered files for actors registered after this commit (incl. this verifier), not on the leaf.

## Links

### Out

- verifies [daemon-resume-once](../capability/rhiza/orchestration/daemon/daemon-resume-once.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
