# api-error-reprompt-plan

source · `01M7A9AA1VK4BF3JVAWBNPX6W8`

## quotes

```json
[]
```

## raw_sha256

268cebdb8a66e4ea18cc935185b7a7c72bbfb0127c970e16e8d163fa636c70a2

## session

54eab379-f795-4280-92c5-01a1d4926df3

## summary

Autonomous plan run of the bug request; no owner was present and no question was put to the arbiter. Research: the incident transcripts (e.g. quality check dbb3e57f) end in Claude Code's synthetic API-error line (isApiErrorMessage, model <synthetic>, text 'Please run /login · API Error: 403 Attention Required! | Cloudflare'), the same line shape the claudecode adapter's UsageLimit already reads. `rhiza session run` already polls each run session's transcript every 30s (watchLimit, internal/api/session_run.go) for the usage limit. The watchdog's agentNudge (20m) and session.Stalled cover a silent session; the incident's sessions sat idle about 18 minutes, under that cutoff. session.message types only into an idle session and waits on a blocked one. Planned: the adapter reads a transient API error; session run re-prompts an autonomous run's session after it with backoff and, when retries run out (or at once in an owner-attended run), records it so the inbox lists an api-error item. The transcript-silence clause stays with the existing watchdog.

## title

Plan: re-prompt run sessions stopped on an API error (autonomous plan run)

## Links

### In

- derived-from [api-error-reprompt-in-session-run](../decision/api-error-reprompt-in-session-run.md)
- derived-from [api-error-retry-policy](../decision/api-error-retry-policy.md)
- derived-from [api-error-inbox-as-idle-detail](../decision/api-error-inbox-as-idle-detail.md)
- derived-from [api-error-silence-left-to-watchdog](../decision/api-error-silence-left-to-watchdog.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
