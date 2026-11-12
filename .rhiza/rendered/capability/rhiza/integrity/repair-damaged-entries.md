# repair-damaged-entries

capability · `01M76QVDF2XHQYY78GE402FTD4` · parent [integrity](../integrity.md)

## acceptance

A testscript (go test -trimpath ./cmd/rhiza -run TestScripts/<name>) reproduces the incident (an uncommitted entry the index has seen is overwritten with another entry's bytes), rhiza repair --dry-run names it and changes nothing, rhiza repair restores the original bytes at 0444 and rhiza render and fold-ok then pass; it also restores a changed committed entry from HEAD, and reports an entry the index never saw as not restorable with exit non-zero; go test -trimpath ./internal/index/... covers reading an entry's stored bytes; describe-golden.txtar is regenerated

## statement

rhiza repair (with --dry-run) finds the entries render refuses and restores each, from HEAD's bytes when the file is tracked at HEAD, else from the worktree's persistent index entries table when its stored bytes parse canonically under that file name's ID, writing through temp file and rename at 0444, deleting nothing, and naming every entry it could not restore and why; rhiza check stays read-only

## Links

### Out

- decided-by [decision-log-guard](../../../decision/decision-log-guard.md)

### In

- verifies [repair-damaged-entries-verified](../../../evidence/repair-damaged-entries-verified.md)
- verifies [repair-damaged-entries-qa-1](../../../evidence/repair-damaged-entries-qa-1.md)
- verifies [repair-damaged-entries-verified-2](../../../evidence/repair-damaged-entries-verified-2.md)
- verifies [repair-damaged-entries-verified-3](../../../evidence/repair-damaged-entries-verified-3.md)
- verifies [pr-465-fix-an-agent-in-a-run](../../../evidence/pr-465-fix-an-agent-in-a-run.md)

## History

- `01M76QJWKEA2M2G87ESG10A47P` Fix An agent in a run overwrote a log entry with a shell command and nothing stopped it, so the plan PR failed CI and the leaves would have been lost. Evidence 2026-11-11: in plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels) the reviser session 59f5535b-caf3-4183-be66-01a1ccc1895f ran 'python3 -m json.tool .rhiza/log/01M76AZHGY605TXF68M27W4SZV/*.json | grep -n v0.1.0'; the glob expanded to two files and json.tool treats the second argument as its output file, so it wrote the change set's open entry (01M76AZHGZ...) over 01M76B7CYMK7R1EKJSF2CHBKF8 (the latest channel's leaves, 22:21:31Z). The broken file was committed in 'Render the plan' and PR #457's lint failed: 'fail log .rhiza/log/01M76AZHGY605TXF68M27W4SZV/01M76B7CYMK7R1EKJSF2CHBKF8.json: entry: invalid: canon: not canonical'. The orchestrator restored the original bytes from the worktree's persistent index (entries table) and pushed 0f5ae263. Wanted: (1) runs' agent sessions get a guard against writing .rhiza/log outside rhiza (a PreToolUse hook that refuses Write/Edit on .rhiza/log/** and warns on shell commands that redirect or write there, or log files made read-only on disk with rhiza writing through a temp file and rename); (2) the commit and render steps refuse to commit a log entry that does not parse canonically or whose id does not match its file name, naming the entry; (3) a rhiza repair command (or rhiza check --fix) that restores a damaged entry from the persistent index when the index still holds its original bytes; (4) agent guidance says to read entries with rhiza show/slice or cat, never tools that can write.
