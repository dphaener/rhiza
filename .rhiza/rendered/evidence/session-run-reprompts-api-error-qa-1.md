# session-run-reprompts-api-error-qa-1

evidence · `01M7ADEKVGSPC9DR60VP23JF26` · actor [claude-code-60ccd5b5](../actor/claude-code-60ccd5b5.md)

## result

PASS

## statement

QA cycle 1: built binary, autonomous run (--autonomous) whose fake claude ended its transcript in the Cloudflare API error and stayed open. session run recorded idle with no detail and sent three re-prompts via rex send at 12:05:47, 12:08:17 and 12:12:17 (1m, 2m, 4m backoff after the 30s poll), each message 'Your last turn ended on an API error: <text>. Carry on where you left off.'; inbox snapshots every 15s listed no api-error and no stalled item during retries; after the fourth backoff (12:16:47) the record held 'stopped on an API error: <text>', the inbox listed it once, and no fourth message went out. An owner-attended run's session sent no message and was recorded at the first poll. go test -trimpath ./internal/api/ -run SessionRunAPIError passes.

## Links

### Out

- verifies [session-run-reprompts-api-error](../capability/rhiza/workflows/workflow-runs/run-resume/session-run-reprompts-api-error.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
