# inbox-lists-api-error-sessions-qa-1

evidence · `01M7ADEKSPY3C2666ZV6SNKN3E` · actor [claude-code-60ccd5b5](../actor/claude-code-60ccd5b5.md)

## result

PASS

## statement

QA cycle 1: with the built binary, an owner-attended run whose session ended on an API error showed in rhiza inbox --json as kind api-error with the session, run, step build, detail 'session <id> stopped on an API error: Please run /login · API Error: 403 ...' and next 'rhiza session message <id> <text>'; the text inbox showed the same; no stalled item. A working hook report (rhiza session hook UserPromptSubmit) cleared the item. Scoped tests for internal/session, internal/tui/inboxapp, internal/api -run 'WebInbox|InboxNotice' and TestScripts/inbox-api-error and describe-golden pass.

## Links

### Out

- verifies [inbox-lists-api-error-sessions](../capability/rhiza/orchestration/inbox/inbox-lists-api-error-sessions.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
