# inbox-nudge-stalled-qa-2

evidence · `01M6DPJRRNDRFKV2P4SFZM81YE` · actor [claude-code-b3f5c4d8](../actor/claude-code-b3f5c4d8.md)

## result

PASS

## statement

At 7f653b81, with the built binary against this live run session in Rex: `rhiza session pane <own session> --json` returned the pane's real text (the Claude Code prompt and status line) with usage_limit false; `rhiza session key <session> escape --dry-run --json` returned dry_run true and sent nothing; key 'tab' was refused ('want escape or enter', exit 2). Escape into a live agent pane and the inbox app's n on a real stalled item were not run (no stalled session to use without disrupting a live agent), so the escape-then-message order rests on go test ./internal/tui/inboxapp, ./internal/api and TestScripts/session-pane-key, all passing.

## Links

### Out

- verifies inbox-nudge-stalled (retired)

## History

- `01M6CRZTNDTPF1SD2WGBY4835Y` An autonomy app edits every setting, starts or stops the daemon, and the inbox nudges stalled sessions
