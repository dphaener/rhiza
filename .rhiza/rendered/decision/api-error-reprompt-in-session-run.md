# api-error-reprompt-in-session-run

decision · `01M7A9C46VS0N9X4DVF0FMYSZE`

## reasoning

The incident transcripts (quality check dbb3e57f) end in Claude Code's synthetic API-error line (isApiErrorMessage, model <synthetic>, text 'Please run /login · API Error: 403 Attention Required! | Cloudflare'), the same shape claudecode's UsageLimit already reads, so the adapter can tell it apart from the model's prose. watchLimit in internal/api/session_run.go already reads every run session's transcript every 30s while its agent runs, attended or autonomous, so the check costs no new read and acts within a minute; the daemon would need a new per-session transcript read on every poll and only runs autonomous orchestration. A usage-limit line is also a synthetic API error and keeps its own pause path, so it is excluded.

## statement

A run session stopped on a transient API error is found and re-prompted by the `rhiza session run` process that started it, from its transcript, on the same poll that already watches for the usage limit, not by the daemon

## Links

### Out

- derived-from [api-error-reprompt-plan](../source/api-error-reprompt-plan.md)

### In

- decided-by [adapter-reads-api-error](../capability/rhiza/workflows/workflow-runs/run-resume/adapter-reads-api-error.md)
- decided-by [session-run-reprompts-api-error](../capability/rhiza/workflows/workflow-runs/run-resume/session-run-reprompts-api-error.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
