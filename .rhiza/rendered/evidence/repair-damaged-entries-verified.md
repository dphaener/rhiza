# repair-damaged-entries-verified

evidence · `01M76Y482KEAVSJN7T1V4BXWW6` · actor [claude-code-f88c07e6](../actor/claude-code-f88c07e6.md)

## result

PASS

## statement

Verified at c2db23b07 (attempt 3). go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(repair|describe-golden)$' passes: the repair script reproduces the incident (an uncommitted, index-seen entry overwritten with another entry's bytes; render refuses with 'damaged log'), repair --dry-run names the entry (path and from=index in --json, 'would restore from index' in text) and leaves bytes and mode 644 unchanged, repair restores the original bytes (cmp) at 0444, then render and fold-ok pass; a changed committed entry is restored from HEAD at 0444; an entry the index never saw exits non-zero with 'not restorable' naming it. go test ./internal/index/... passes (TestStoredEntry covers reading stored bytes). The earlier finding is fixed: go test ./internal/api -run TestEveryResponseIsHumaner passes, RepairResponse now has Human(io.Writer). Full ./internal/api tests all PASS; its root-state guard flagged ~/code/.rhiza-root/usage.live.json changing during the 9-minute run, which the live root daemon writes and repair does not touch. Read the diff of internal/api/repair.go: HEAD bytes first, else index bytes only when they decode canonically with the path's change set and ID, written via store.WriteEntryFile, nothing deleted.

## Links

### Out

- verifies [repair-damaged-entries](../capability/rhiza/integrity/repair-damaged-entries.md)

## History

- `01M76QJWKEA2M2G87ESG10A47P` Fix An agent in a run overwrote a log entry with a shell command and nothing stopped it, so the plan PR failed CI and the leaves would have been lost. Evidence 2026-11-11: in plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels) the reviser session 59f5535b-caf3-4183-be66-01a1ccc1895f ran 'python3 -m json.tool .rhiza/log/01M76AZHGY605TXF68M27W4SZV/*.json | grep -n v0.1.0'; the glob expanded to two files and json.tool treats the second argument as its output file, so it wrote the change set's open entry (01M76AZHGZ...) over 01M76B7CYMK7R1EKJSF2CHBKF8 (the latest channel's leaves, 22:21:31Z). The broken file was committed in 'Render the plan' and PR #457's lint failed: 'fail log .rhiza/log/01M76AZHGY605TXF68M27W4SZV/01M76B7CYMK7R1EKJSF2CHBKF8.json: entry: invalid: canon: not canonical'. The orchestrator restored the original bytes from the worktree's persistent index (entries table) and pushed 0f5ae263. Wanted: (1) runs' agent sessions get a guard against writing .rhiza/log outside rhiza (a PreToolUse hook that refuses Write/Edit on .rhiza/log/** and warns on shell commands that redirect or write there, or log files made read-only on disk with rhiza writing through a temp file and rename); (2) the commit and render steps refuse to commit a log entry that does not parse canonically or whose id does not match its file name, naming the entry; (3) a rhiza repair command (or rhiza check --fix) that restores a damaged entry from the persistent index when the index still holds its original bytes; (4) agent guidance says to read entries with rhiza show/slice or cat, never tools that can write.
