# api-error-silence-left-to-watchdog

decision · `01M7A9C4A6P40CDWVEP31ZJXQB`

## reasoning

internal/api/watchdog_policy.go DefaultAgentNudge and internal/session/status.go Stalled (from change set 01M6N9AF69) flag an idle or working session silent past 20 minutes in an autonomous run; the incident's sessions sat about 18 minutes, under that cutoff, before the orchestrator stepped in. The missing piece is recognising the API-error stop at once, which needs no silence threshold or process inspection.

## statement

The intent's other trigger, a transcript that has not grown for N minutes with no child process, is not built here: the watchdog's agentNudge (20m by default) and session.Stalled already nudge and list a silent autonomous run session; this change set catches the API-error stop itself

## Links

### Out

- derived-from [api-error-reprompt-plan](../source/api-error-reprompt-plan.md)

### In

- decided-by [inbox-lists-api-error-sessions](../capability/rhiza/orchestration/inbox/inbox-lists-api-error-sessions.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
