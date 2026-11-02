# inbox-nudge-stalled-verified-2

evidence · `01M6DPACFND8M7V0SN9PYG3TQ9` · actor [claude-code-53813c58](../actor/claude-code-53813c58.md)

## result

PASS

## statement

At 2581003e, go test -count=1 ./internal/tui/inboxapp -run Nudge passed TestNudgeSendsCarryOn (n on a stalled item: session.pane then session.message with the carry-on text to SESSN1), TestNudgeEscapesUsageLimitPromptFirst (usage_limit pane: session.pane, session.key escape, then session.message) and TestNudgeOnlyOnStalled (n on a failed item makes no call), all on a fake runner. The fix for finding f1 also passes: internal/api TestSessionPane*/TestSessionKey* (rex capture/send-key argv and dir, PaneLimiter via the adapter, other keys refused with nothing sent, dry run sends nothing, exited/tabless/unknown refused) and TestScripts/session-pane-key against a fake rex. make fmt-check vet test rc=0; self-check fails only on stale .rhiza/rendered actor files from run model writes. Read the diff of the leaf's commits 2f4d0eb2 and 2581003e.

## Links

### Out

- verifies [inbox-nudge-stalled](../capability/rhiza/orchestration/inbox/inbox-list/inbox-stalled-items/inbox-nudge-stalled.md)

## History

- `01M6CRZTNDTPF1SD2WGBY4835Y` An autonomy app edits every setting, starts or stops the daemon, and the inbox nudges stalled sessions
