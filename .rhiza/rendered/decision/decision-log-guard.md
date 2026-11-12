# decision-log-guard

decision · `01M76QSWXXZ9AZCW4918RBPYYH`

## reasoning

The incident was a shell command (python3 -m json.tool with a glob), so a PreToolUse hook on Write/Edit alone would not have stopped it, and parsing shell for writes is unreliable; read-only entries stop any non-root writer, and rhiza never rewrites an entry in place (writeNew links a temp file and refuses an existing name; removal and git checkout unlink, which file mode does not block), so 0444 costs nothing. The overwritten entry had been written by rhiza in the same run, so 0444 at write time would have refused it; render re-applying the mode covers checked-out entries. The question notes every commit step runs rhiza render first and render already reads the log but ignores its problems, so making render fail is the single choke point the change set's intent (2) asks for, naming the entry. For repair, the orchestrator's manual fix was restoring from the persistent index entries table (index catchup stores only good entries, so it keeps original bytes); HEAD is the authority for committed entries. A separate rhiza repair keeps rhiza check a read-only CI lint, and --dry-run lets a run see what would change. Guidance (4) is the intent's own ask.

## statement

(1) Guard in two layers. rhiza writes every log entry read-only: both entry writers set the temp file to 0444 before putting it in place, through one shared helper: store's writeNew and internal/api/session.go's copyNewFile (which copies a change set's entries into a session's worktree, so a copied entry is read-only while the agent session runs), and rhiza render re-applies 0444 to every entry under .rhiza/log (git checkouts give 0644). That blocks shell writes (redirects, json.tool's output argument, sed -i on a path) without parsing shell. In addition the claudecode adapter's PreToolUse hook (rhiza session hook) denies Write, Edit, MultiEdit and NotebookEdit on any path under .rhiza/log/, with a reason that says to write through rhiza apply and read with rhiza show/slice or cat. No shell-command parsing or warning: the file mode covers shell. (2) rhiza render refuses: it already reads the whole log, so it fails (non-zero exit, nothing rendered) when any entry does not parse canonically, its id does not match its file name, or a log file tracked at HEAD differs from HEAD, naming each entry's path and its problem. Every run commit step runs render first, so this one check covers build/review/ship-change render, the plan gate's 'Render the plan', ship land and ship rebase. (3) A new command, rhiza repair (with --dry-run), finds the same damaged entries and restores each: a file tracked at HEAD from HEAD's bytes; otherwise from the persistent index's entries table when its stored raw parses canonically and its id matches the file name. It writes through temp file + rename at 0444, never deletes, and reports every entry it could not restore and why. rhiza check stays read-only. (4) The rhiza guide and the agent skills say to read log entries with rhiza show/slice or cat, never with tools that can write (json.tool, jq -i, sed -i, redirects).

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-947b0dc6-1](../question/question-947b0dc6-1.md)
- derived-from [log-guard-plan](../source/log-guard-plan.md)

### In

- decided-by [log-entries-read-only](../capability/rhiza/model/operation-log/log-entries-read-only.md)
- decided-by [session-hook-denies-log-edits](../capability/rhiza/runtime/rex-sessions/session-hook-denies-log-edits.md)
- decided-by [render-refuses-damaged-log](../capability/rhiza/integrity/render/render-refuses-damaged-log.md)
- decided-by [repair-damaged-entries](../capability/rhiza/integrity/repair-damaged-entries.md)
- decided-by [guide-reads-log-safely](../capability/rhiza/guide/agent-guide/guide-reads-log-safely.md)

## History

- `01M76QJWKEA2M2G87ESG10A47P` Fix An agent in a run overwrote a log entry with a shell command and nothing stopped it, so the plan PR failed CI and the leaves would have been lost. Evidence 2026-11-11: in plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels) the reviser session 59f5535b-caf3-4183-be66-01a1ccc1895f ran 'python3 -m json.tool .rhiza/log/01M76AZHGY605TXF68M27W4SZV/*.json | grep -n v0.1.0'; the glob expanded to two files and json.tool treats the second argument as its output file, so it wrote the change set's open entry (01M76AZHGZ...) over 01M76B7CYMK7R1EKJSF2CHBKF8 (the latest channel's leaves, 22:21:31Z). The broken file was committed in 'Render the plan' and PR #457's lint failed: 'fail log .rhiza/log/01M76AZHGY605TXF68M27W4SZV/01M76B7CYMK7R1EKJSF2CHBKF8.json: entry: invalid: canon: not canonical'. The orchestrator restored the original bytes from the worktree's persistent index (entries table) and pushed 0f5ae263. Wanted: (1) runs' agent sessions get a guard against writing .rhiza/log outside rhiza (a PreToolUse hook that refuses Write/Edit on .rhiza/log/** and warns on shell commands that redirect or write there, or log files made read-only on disk with rhiza writing through a temp file and rename); (2) the commit and render steps refuse to commit a log entry that does not parse canonically or whose id does not match its file name, naming the entry; (3) a rhiza repair command (or rhiza check --fix) that restores a damaged entry from the persistent index when the index still holds its original bytes; (4) agent guidance says to read entries with rhiza show/slice or cat, never tools that can write.
