# session-hook-denies-log-edits-qa-1

evidence · `01M76YJF80QBFE0NDE9ER78FAS` · actor [claude-code-64f8f174](../actor/claude-code-64f8f174.md)

## result

PASS

## statement

QA at 044b0aef3, binary built with go build -trimpath into a scratch dir, run by hand as uid 1000 in a scratch git repo (RHIZA_* unset): In a scratch session (rhiza session start --permission-mode manual --step solo, fake rex), fed rhiza session hook PreToolUse payloads: Write on .rhiza/log/a/b.json, Edit on the absolute worktree path, MultiEdit on ../x/.rhiza/log/a/b.json and NotebookEdit (notebook_path) under .rhiza/log all got permissionDecision deny with the reason 'Never edit .rhiza/log by hand ... Write through `rhiza apply`; read with `rhiza show`, `rhiza slice` or `cat`.' Write on .rhiza/logx/a.json and notes.md, Read of a log entry and Bash 'echo > .rhiza/log/...' got no reply. go test ./internal/adapter/... and TestScripts/session-hook-log-guard ok.

## Links

### Out

- verifies [session-hook-denies-log-edits](../capability/rhiza/runtime/rex-sessions/session-hook-denies-log-edits.md)

## History

- `01M76QJWKEA2M2G87ESG10A47P` Fix An agent in a run overwrote a log entry with a shell command and nothing stopped it, so the plan PR failed CI and the leaves would have been lost. Evidence 2026-11-11: in plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels) the reviser session 59f5535b-caf3-4183-be66-01a1ccc1895f ran 'python3 -m json.tool .rhiza/log/01M76AZHGY605TXF68M27W4SZV/*.json | grep -n v0.1.0'; the glob expanded to two files and json.tool treats the second argument as its output file, so it wrote the change set's open entry (01M76AZHGZ...) over 01M76B7CYMK7R1EKJSF2CHBKF8 (the latest channel's leaves, 22:21:31Z). The broken file was committed in 'Render the plan' and PR #457's lint failed: 'fail log .rhiza/log/01M76AZHGY605TXF68M27W4SZV/01M76B7CYMK7R1EKJSF2CHBKF8.json: entry: invalid: canon: not canonical'. The orchestrator restored the original bytes from the worktree's persistent index (entries table) and pushed 0f5ae263. Wanted: (1) runs' agent sessions get a guard against writing .rhiza/log outside rhiza (a PreToolUse hook that refuses Write/Edit on .rhiza/log/** and warns on shell commands that redirect or write there, or log files made read-only on disk with rhiza writing through a temp file and rename); (2) the commit and render steps refuse to commit a log entry that does not parse canonically or whose id does not match its file name, naming the entry; (3) a rhiza repair command (or rhiza check --fix) that restores a damaged entry from the persistent index when the index still holds its original bytes; (4) agent guidance says to read entries with rhiza show/slice or cat, never tools that can write.
