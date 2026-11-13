# api-error-retry-policy

decision · `01M7A9C47YW5N5EHZZHGEHD1MK`

## reasoning

The intent asks for 'retrying a few times with backoff, then raises an inbox item'. The incident's 403s were transient: each session resumed within seconds of a message about 15 minutes later, so short waits doubling from 1 minute cover a passing outage in about 7 minutes without hammering a failing API. A real outage or an expired login ('Please run /login') outlasts three tries and then needs a person. Keying the count on consecutive errors means a session that recovered and later hits a new error gets fresh retries. Re-prompting only autonomous runs follows idle-stall-autonomous-only and daemon-nudges-stalled: an owner-attended session may be the owner's own to resume, so it is surfaced, not typed into.

## statement

In a session an autonomous run launched, an API-error stop is re-prompted at most 3 times in a row, after waits of 1, 2 and 4 minutes from when the error is seen; an assistant line that is not an API error resets the count; a fourth consecutive API error, or any API error in a session of an owner-attended run, is not re-prompted and goes to the inbox at once

## Links

### Out

- derived-from [api-error-reprompt-plan](../source/api-error-reprompt-plan.md)

### In

- decided-by [session-run-reprompts-api-error](../capability/rhiza/workflows/workflow-runs/run-resume/session-run-reprompts-api-error.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
