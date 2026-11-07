# usage-live-snapshot-record

capability · `01M6T2P2NB2AXZ8SE7K6XJMK1C` · parent [usage-governor](../usage-governor.md)

## acceptance

internal/adapter/claudecode tests show Settings carries a statusLine command alongside the hooks; a new usage-live-record.txtar pipes status line JSON with five_hour and seven_day rate_limits into the command and finds one latest snapshot per kind in the root's state folder (and the repo's daemon dir outside a root), a second input replaces the first, and input without rate_limits or malformed JSON exits 0 and records nothing; with a fake user settings file whose statusLine command echoes its stdin the command prints that output after recording, with no user statusLine it prints nothing, and a failing user command still exits 0 (go test -trimpath ./internal/adapter/claudecode/... and go test -trimpath ./cmd/rhiza -run TestScripts/usage-live-record)

## statement

Every Claude Code session Rhiza launches runs a rhiza statusLine command (added to the adapter's --settings JSON) that reads the status line JSON on stdin and records each rate_limits window it carries (five_hour as session, seven_day as weekly: used_percentage, resets_at and the capture time) to a live-usage file in the same state folder as the usage-limit hits, keeping the latest snapshot per window kind; input with no rate_limits records nothing. Since the --settings statusLine replaces the user's own, the command then runs the statusLine command the user's own Claude Code settings configure, if any, with the same stdin, and prints its output as the status line; with none it prints nothing. The command never fails the session, even when the user's command fails

## Links

### Out

- decided-by [live-usage-source](../../../../decision/live-usage-source.md)

### In

- depends-on [usage-live-snapshot-read](usage-live-snapshot-read.md)
- verifies [usage-live-snapshot-record-verified](../../../../evidence/usage-live-snapshot-record-verified.md)
- verifies [usage-live-snapshot-record-qa-1](../../../../evidence/usage-live-snapshot-record-qa-1.md)
- verifies [pr-349-fix-the-usage-governor-s](../../../../evidence/pr-349-fix-the-usage-governor-s.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
