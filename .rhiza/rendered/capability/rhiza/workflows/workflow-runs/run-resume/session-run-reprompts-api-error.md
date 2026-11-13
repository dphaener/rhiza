# session-run-reprompts-api-error

capability · `01M7A9C4DABT4TTGNKQ5JFRSV9` · parent [run-resume](../run-resume.md)

## acceptance

A Go test in internal/api/session_run_apierror_test.go (go test -trimpath ./internal/api/ -run SessionRunAPIError) with a fake adapter, a fake process, a fake message sender and the poll and backoff shortened: an autonomous session whose new transcript ends in an API error gets one message naming the error after the first backoff and none before it; its record, last reported working by the hook, holds an idle hook report with no 'stopped on an API error: ' detail before each send, and inboxItems over that run lists no api-error item for the session while retries remain; an assistant line written during the backoff cancels the send; three consecutive API errors give three messages and a fourth gives none, with the record left on 'stopped on an API error: ' and the error text; an error, a recovered assistant line, then a new error re-prompts again; an owner-attended session's API error gives no message and the record at once; a usage-limit line gives no message and the usage-limit pause still happens; an API error in the transcript before this start's base gives nothing. Touches internal/api/session_run.go, internal/api/session_run_apierror.go, internal/api/session_run_apierror_test.go

## statement

While its agent runs, `rhiza session run` watches the session's transcript past where this start began, on the usage-limit poll; when it ends in an API error (adapter-reads-api-error) in a session an autonomous run launched, it waits the retry policy's backoff, then, if no later assistant line has appeared, records a plain idle hook report, with no 'stopped on an API error: ' detail so the inbox lists nothing while retries remain, unless the status already resolves idle and sends the session, through session.message, a message naming the error and telling it to carry on where it left off; it re-prompts at most 3 consecutive errors, and on the fourth, or on the first in an owner-attended run's session, it re-prompts nothing and records the session for the inbox (inbox-lists-api-error-sessions)

## Links

### Out

- depends-on [adapter-reads-api-error](adapter-reads-api-error.md)
- depends-on [inbox-lists-api-error-sessions](../../../orchestration/inbox/inbox-lists-api-error-sessions.md)
- decided-by [api-error-reprompt-in-session-run](../../../../../decision/api-error-reprompt-in-session-run.md)
- decided-by [api-error-retry-policy](../../../../../decision/api-error-retry-policy.md)

### In

- verifies [session-run-reprompts-api-error-verified](../../../../../evidence/session-run-reprompts-api-error-verified.md)
- verifies [session-run-reprompts-api-error-qa-1](../../../../../evidence/session-run-reprompts-api-error-qa-1.md)
- verifies [session-run-reprompts-api-error-verified-2](../../../../../evidence/session-run-reprompts-api-error-verified-2.md)
- verifies [pr-548-fix-run-sessions-that-hit-an](../../../../../evidence/pr-548-fix-run-sessions-that-hit-an.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
