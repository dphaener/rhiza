# inbox-nudge-stalled-qa-1

evidence · `01M6DMT5AP9NFMDDQ2MTHCGQJ9` · actor [claude-code-f364d51f](../actor/claude-code-f364d51f.md)

## result

PASS

## statement

At 69611388, with the built binary: `rhiza session key ZZZZ tab` refuses with 'want escape or enter' (exit 2); `rhiza session pane ZZZZ` says no session matches (exit 4); `rhiza describe session.key` documents escape/enter. Reading or pressing keys in a live agent pane was not run (driving a live session's terminal was not permitted in this QA session), so the pane/escape/message order rests on go test ./internal/tui/inboxapp (pass) and TestPaneUsageLimit (pass) plus reading the diff: n acts only on stalled items, calls session.pane, session.key escape when usage_limit, then session.message.

## Links

### Out

- verifies inbox-nudge-stalled (retired)

## History

- `01M6CRZTNDTPF1SD2WGBY4835Y` An autonomy app edits every setting, starts or stops the daemon, and the inbox nudges stalled sessions
