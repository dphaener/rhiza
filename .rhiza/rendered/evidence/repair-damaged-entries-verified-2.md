# repair-damaged-entries-verified-2

evidence · `01M76ZYVSFTX7MYDE8SVP9DBHB` · actor [claude-code-02e43888](../actor/claude-code-02e43888.md)

## result

PASS

## statement

Verified at 4da6a0050 (fix for quality f1: HEAD bytes restored unchecked). Read the diff: originalBytes now applies isEntryBytes (DecodeCanonical, ID and change set match the path) to HEAD's blob for entry paths, falls back to the index, and reports 'HEAD's copy is damaged too and ...' when both fail; non-entry files such as repo.json still take HEAD unchecked. go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(repair|describe-golden|render-refuses-damaged-log)$': ok. The repair script covers the incident (an uncommitted entry overwritten with another entry's bytes): dry-run names it and changes nothing, repair restores it from the index at 0444, then render and fold-ok pass. A changed committed entry is restored from HEAD. New case: an entry committed already damaged makes render refuse; repair says 'restored from index', not 'from HEAD', and cmp matches the original bytes at 0444; render and fold-ok pass after commit. An entry the index never saw is reported not restorable, exit non-zero. go test ./internal/index/...: ok. go test ./internal/api -run 'Repair|Humaner|Describe': ok. gofmt clean, go vet ./internal/api clean. The whole internal/api package hit a 240s timeout on a loaded shared host (load 5.9 on 4 cores) while TestUsageWindowsAnchored (unrelated) had run for 2s; no test failed, so CI's make check covers it.

## Links

### Out

- verifies [repair-damaged-entries](../capability/rhiza/integrity/repair-damaged-entries.md)

## History

- `01M76QJWKEA2M2G87ESG10A47P` Fix An agent in a run overwrote a log entry with a shell command and nothing stopped it, so the plan PR failed CI and the leaves would have been lost. Evidence 2026-11-11: in plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels) the reviser session 59f5535b-caf3-4183-be66-01a1ccc1895f ran 'python3 -m json.tool .rhiza/log/01M76AZHGY605TXF68M27W4SZV/*.json | grep -n v0.1.0'; the glob expanded to two files and json.tool treats the second argument as its output file, so it wrote the change set's open entry (01M76AZHGZ...) over 01M76B7CYMK7R1EKJSF2CHBKF8 (the latest channel's leaves, 22:21:31Z). The broken file was committed in 'Render the plan' and PR #457's lint failed: 'fail log .rhiza/log/01M76AZHGY605TXF68M27W4SZV/01M76B7CYMK7R1EKJSF2CHBKF8.json: entry: invalid: canon: not canonical'. The orchestrator restored the original bytes from the worktree's persistent index (entries table) and pushed 0f5ae263. Wanted: (1) runs' agent sessions get a guard against writing .rhiza/log outside rhiza (a PreToolUse hook that refuses Write/Edit on .rhiza/log/** and warns on shell commands that redirect or write there, or log files made read-only on disk with rhiza writing through a temp file and rename); (2) the commit and render steps refuse to commit a log entry that does not parse canonically or whose id does not match its file name, naming the entry; (3) a rhiza repair command (or rhiza check --fix) that restores a damaged entry from the persistent index when the index still holds its original bytes; (4) agent guidance says to read entries with rhiza show/slice or cat, never tools that can write.
