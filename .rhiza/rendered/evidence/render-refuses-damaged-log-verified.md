# render-refuses-damaged-log-verified

evidence · `01M76VDMGNKZ4G1MART72QN5EK` · actor [claude-code-1a1d9469](../actor/claude-code-1a1d9469.md)

## result

PASS

## statement

Verifier at 7742d897c: go test -trimpath ./cmd/rhiza -run TestScripts/(render-refuses-damaged-log|log-entries-read-only|render)$ passes; the new script makes all three damages in turn, asserts render exits non-zero with the problem, diffs .rhiza/rendered unchanged each time, and a restored log renders again. By hand with a scratch build: on committed entries, a non-canonical rewrite and one entry's bytes under another's name each fail render naming the path with 'not canonical' / 'id or change set does not match its path' (plus 'differs from HEAD'), and a canonical sed edit fails with 'differs from HEAD'; on uncommitted entries (no HEAD diff) the non-canonical and ID-mismatch damages still fail with path and problem, so each check stands alone; .rhiza/rendered stayed unchanged and a clean log rendered (exit 0). internal/store passes; internal/api render/log/read-only tests (39) pass. The full internal/api package did not finish on the overloaded host (load 18 on 4 cores; TestReadBudget timed out), and rootguard flags the known usage.live.json leak into ~/code/.rhiza-root, unrelated to this diff. Nit: the script's cases 1 and 2 assert the problem text but not the path; case 3 asserts the path.

## Links

### Out

- verifies [render-refuses-damaged-log](../capability/rhiza/integrity/render/render-refuses-damaged-log.md)

## History

- `01M76QJWKEA2M2G87ESG10A47P` Fix An agent in a run overwrote a log entry with a shell command and nothing stopped it, so the plan PR failed CI and the leaves would have been lost. Evidence 2026-11-11: in plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels) the reviser session 59f5535b-caf3-4183-be66-01a1ccc1895f ran 'python3 -m json.tool .rhiza/log/01M76AZHGY605TXF68M27W4SZV/*.json | grep -n v0.1.0'; the glob expanded to two files and json.tool treats the second argument as its output file, so it wrote the change set's open entry (01M76AZHGZ...) over 01M76B7CYMK7R1EKJSF2CHBKF8 (the latest channel's leaves, 22:21:31Z). The broken file was committed in 'Render the plan' and PR #457's lint failed: 'fail log .rhiza/log/01M76AZHGY605TXF68M27W4SZV/01M76B7CYMK7R1EKJSF2CHBKF8.json: entry: invalid: canon: not canonical'. The orchestrator restored the original bytes from the worktree's persistent index (entries table) and pushed 0f5ae263. Wanted: (1) runs' agent sessions get a guard against writing .rhiza/log outside rhiza (a PreToolUse hook that refuses Write/Edit on .rhiza/log/** and warns on shell commands that redirect or write there, or log files made read-only on disk with rhiza writing through a temp file and rename); (2) the commit and render steps refuse to commit a log entry that does not parse canonically or whose id does not match its file name, naming the entry; (3) a rhiza repair command (or rhiza check --fix) that restores a damaged entry from the persistent index when the index still holds its original bytes; (4) agent guidance says to read entries with rhiza show/slice or cat, never tools that can write.
