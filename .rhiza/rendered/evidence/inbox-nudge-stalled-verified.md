# inbox-nudge-stalled-verified

evidence · `01M6DMJ4P62YBA9TDTJJ4BYES1` · actor [claude-code-7bfccb31](../actor/claude-code-7bfccb31.md)

## result

PASS

## statement

At 2f4d0eb2, go test ./internal/tui/inboxapp -run Nudge -v passes TestNudgeSendsCarryOn (on the stalled item, n calls session.pane then session.message SESSN1 with the carry-on text), TestNudgeEscapesUsageLimitPromptFirst (pane with usage_limit true: session.pane, session.key escape, then session.message) and TestNudgeOnlyOnStalled (n on a failed item makes no calls), all on a fake runner. TestPaneUsageLimit in claudecode passes (only the last 15 pane lines count). gofmt is clean and go test ./... passes. Read the diff: n is handled only for api.InboxStalled, session.pane captures the pane with rex capture and asks the adapter through the optional PaneLimiter, session.key accepts only escape or enter and sends it with rex send-key. make self-check fails on a stale .rhiza/rendered for actor claude-code-daace8aa, the same failure as at parent 5318fc77, so this leaf did not cause it.

## Links

### Out

- verifies inbox-nudge-stalled (retired)

## History

- `01M6CRZTNDTPF1SD2WGBY4835Y` An autonomy app edits every setting, starts or stops the daemon, and the inbox nudges stalled sessions
