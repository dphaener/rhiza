# log-entries-read-only

capability · `01M76QVD40AW1RFPN3ENSY5F9P` · parent [operation-log](../operation-log.md)

## acceptance

go test -trimpath ./internal/store/... has a test that a written entry is 0444 and that opening it for writing fails (as non-root); go test -trimpath ./internal/api/... has a test that an entry copied by copyChangeSetLog into a session worktree is 0444; a testscript in cmd/rhiza/testdata/script (go test -trimpath ./cmd/rhiza -run TestScripts/<name>) chmods an entry to 0644, runs rhiza render and finds it 0444 again

## statement

rhiza writes every log entry file read-only (0444): store's writeNew and the session worktree copy (internal/api/session.go copyNewFile, via copyChangeSetLog) both write entries through one shared helper that sets 0444; rhiza render sets 0444 on every entry under .rhiza/log, so a shell command that opens an entry for writing fails instead of overwriting it

## Links

### Out

- decided-by [decision-log-guard](../../../../decision/decision-log-guard.md)

### In

- verifies [log-entries-read-only-verified](../../../../evidence/log-entries-read-only-verified.md)
- verifies [log-entries-read-only-qa-1](../../../../evidence/log-entries-read-only-qa-1.md)
- verifies [pr-465-fix-an-agent-in-a-run](../../../../evidence/pr-465-fix-an-agent-in-a-run.md)

## History

- `01M76QJWKEA2M2G87ESG10A47P` Fix An agent in a run overwrote a log entry with a shell command and nothing stopped it, so the plan PR failed CI and the leaves would have been lost. Evidence 2026-11-11: in plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels) the reviser session 59f5535b-caf3-4183-be66-01a1ccc1895f ran 'python3 -m json.tool .rhiza/log/01M76AZHGY605TXF68M27W4SZV/*.json | grep -n v0.1.0'; the glob expanded to two files and json.tool treats the second argument as its output file, so it wrote the change set's open entry (01M76AZHGZ...) over 01M76B7CYMK7R1EKJSF2CHBKF8 (the latest channel's leaves, 22:21:31Z). The broken file was committed in 'Render the plan' and PR #457's lint failed: 'fail log .rhiza/log/01M76AZHGY605TXF68M27W4SZV/01M76B7CYMK7R1EKJSF2CHBKF8.json: entry: invalid: canon: not canonical'. The orchestrator restored the original bytes from the worktree's persistent index (entries table) and pushed 0f5ae263. Wanted: (1) runs' agent sessions get a guard against writing .rhiza/log outside rhiza (a PreToolUse hook that refuses Write/Edit on .rhiza/log/** and warns on shell commands that redirect or write there, or log files made read-only on disk with rhiza writing through a temp file and rename); (2) the commit and render steps refuse to commit a log entry that does not parse canonically or whose id does not match its file name, naming the entry; (3) a rhiza repair command (or rhiza check --fix) that restores a damaged entry from the persistent index when the index still holds its original bytes; (4) agent guidance says to read entries with rhiza show/slice or cat, never tools that can write.
