# log-entries-read-only-verified

evidence · `01M76VG6SHJDRZJM0X8ZYJP98Q` · actor [claude-code-6339fce5](../actor/claude-code-6339fce5.md)

## result

PASS

## statement

Verifier claude-code-6339fce5 at 5396151294: read the diff (store.writeNew and the new store.WriteEntryFile share writeTemp with EntryMode 0444, repo.json kept 0644; session.go copyNewFile now calls store.WriteEntryFile; runRender calls store.SetEntryModes on .rhiza/log). As uid 1000: go test -trimpath ./internal/store/... ok (TestWrite asserts mode 444 and that os.OpenFile O_WRONLY|O_TRUNC fails); go test -trimpath ./internal/api -run TestCopyChangeSetLogWritesReadOnly PASS (copyChangeSetLog output is 0444); go test -trimpath ./cmd/rhiza -run TestScripts/log-entries-read-only PASS (entries 444 after add --record, chmod 644 then rhiza render restores 444). go test -trimpath -short ./internal/api ran every test with no failures; the only red was the package's root-state guard reporting ~/code/.rhiza-root/usage.live.json changed during the 8-minute run, live daemon usage state this diff does not touch. The non-short run stalled in TestInboxGateQuestionsBudget's quiet-machine wait (load ~7 on 4 cores), unrelated to the change.

## Links

### Out

- verifies [log-entries-read-only](../capability/rhiza/model/operation-log/log-entries-read-only.md)

## History

- `01M76QJWKEA2M2G87ESG10A47P` Fix An agent in a run overwrote a log entry with a shell command and nothing stopped it, so the plan PR failed CI and the leaves would have been lost. Evidence 2026-11-11: in plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels) the reviser session 59f5535b-caf3-4183-be66-01a1ccc1895f ran 'python3 -m json.tool .rhiza/log/01M76AZHGY605TXF68M27W4SZV/*.json | grep -n v0.1.0'; the glob expanded to two files and json.tool treats the second argument as its output file, so it wrote the change set's open entry (01M76AZHGZ...) over 01M76B7CYMK7R1EKJSF2CHBKF8 (the latest channel's leaves, 22:21:31Z). The broken file was committed in 'Render the plan' and PR #457's lint failed: 'fail log .rhiza/log/01M76AZHGY605TXF68M27W4SZV/01M76B7CYMK7R1EKJSF2CHBKF8.json: entry: invalid: canon: not canonical'. The orchestrator restored the original bytes from the worktree's persistent index (entries table) and pushed 0f5ae263. Wanted: (1) runs' agent sessions get a guard against writing .rhiza/log outside rhiza (a PreToolUse hook that refuses Write/Edit on .rhiza/log/** and warns on shell commands that redirect or write there, or log files made read-only on disk with rhiza writing through a temp file and rename); (2) the commit and render steps refuse to commit a log entry that does not parse canonically or whose id does not match its file name, naming the entry; (3) a rhiza repair command (or rhiza check --fix) that restores a damaged entry from the persistent index when the index still holds its original bytes; (4) agent guidance says to read entries with rhiza show/slice or cat, never tools that can write.
