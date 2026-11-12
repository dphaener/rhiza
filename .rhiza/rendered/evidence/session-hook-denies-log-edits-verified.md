# session-hook-denies-log-edits-verified

evidence · `01M76VYRW6ZQJKVS4M9EX3H2C5` · actor [claude-code-e738d61c](../actor/claude-code-e738d61c.md)

## result

PASS

## statement

Verifier claude-code-e738d61c at 2e51e76c6: go test -trimpath -count=1 ./internal/adapter/... passes; TestGuardReply denies Write (absolute and relative), Edit (absolute), MultiEdit (relative ../), NotebookEdit (relative notebook_path) under .rhiza/log/, and gives no reply for Write on /w/.rhiza/model.json, for Bash on a log path, for Read of a log path, and for PostToolUse. go test -trimpath ./cmd/rhiza -run TestScripts/session-hook-log-guard passes: rhiza session hook PreToolUse with a Write on .rhiza/log/x/y.json prints permissionDecision deny naming rhiza apply; Write elsewhere, Bash and Read print nothing. The session-related internal/api tests (-run Session|Hook|Ask|Deliver|Interview) pass. The full internal/api run failed only its guard against the enclosing root's .rhiza-root/usage.live.json changing, which the live root daemon writes; this change does not touch usage. Diff read: session_run.go adds logEditReason (says write through rhiza apply, read with rhiza show, rhiza slice or cat) and sets the hook reply from Adapter.GuardReply; claudecode GuardReply checks only the four file tools' file_path/notebook_path, cleans and resolves the path against cwd, and parses no shell commands.

## Links

### Out

- verifies [session-hook-denies-log-edits](../capability/rhiza/runtime/rex-sessions/session-hook-denies-log-edits.md)

## History

- `01M76QJWKEA2M2G87ESG10A47P` Fix An agent in a run overwrote a log entry with a shell command and nothing stopped it, so the plan PR failed CI and the leaves would have been lost. Evidence 2026-11-11: in plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels) the reviser session 59f5535b-caf3-4183-be66-01a1ccc1895f ran 'python3 -m json.tool .rhiza/log/01M76AZHGY605TXF68M27W4SZV/*.json | grep -n v0.1.0'; the glob expanded to two files and json.tool treats the second argument as its output file, so it wrote the change set's open entry (01M76AZHGZ...) over 01M76B7CYMK7R1EKJSF2CHBKF8 (the latest channel's leaves, 22:21:31Z). The broken file was committed in 'Render the plan' and PR #457's lint failed: 'fail log .rhiza/log/01M76AZHGY605TXF68M27W4SZV/01M76B7CYMK7R1EKJSF2CHBKF8.json: entry: invalid: canon: not canonical'. The orchestrator restored the original bytes from the worktree's persistent index (entries table) and pushed 0f5ae263. Wanted: (1) runs' agent sessions get a guard against writing .rhiza/log outside rhiza (a PreToolUse hook that refuses Write/Edit on .rhiza/log/** and warns on shell commands that redirect or write there, or log files made read-only on disk with rhiza writing through a temp file and rename); (2) the commit and render steps refuse to commit a log entry that does not parse canonically or whose id does not match its file name, naming the entry; (3) a rhiza repair command (or rhiza check --fix) that restores a damaged entry from the persistent index when the index still holds its original bytes; (4) agent guidance says to read entries with rhiza show/slice or cat, never tools that can write.
