# pr-548-fix-run-sessions-that-hit-an

evidence · `01M7ANE3WAZ4WWGQJDFJBD0WVQ` · actor [dphaener](../actor/dphaener.md)

## head

e6c3b56837e808621ada348c94f7acf4e105fede

## pr

https://github.com/dphaener/rhiza/pull/548

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #548 (https://github.com/dphaener/rhiza/pull/548) passed CI at e6c3b56

## status

verified

## Links

### Out

- verifies [adapter-reads-api-error](../capability/rhiza/workflows/workflow-runs/run-resume/adapter-reads-api-error.md)
- verifies [session-run-reprompts-api-error](../capability/rhiza/workflows/workflow-runs/run-resume/session-run-reprompts-api-error.md)
- verifies [inbox-lists-api-error-sessions](../capability/rhiza/orchestration/inbox/inbox-lists-api-error-sessions.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
