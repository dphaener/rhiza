# adapter-reads-api-error-qa-1

evidence · `01M7ADEKQCBSNDR8473DKAV91R` · actor [claude-code-60ccd5b5](../actor/claude-code-60ccd5b5.md)

## result

PASS

## statement

QA cycle 1: built rhiza from the branch to a scratch path and ran rhiza run start --foreground in a hermetic Rhiza root with a fake rex and a fake claude whose transcript ended in the synthetic isApiErrorMessage line 'Please run /login · API Error: 403 Attention Required! | Cloudflare' after a model line. The Claude Code adapter recognised it end to end: the error text reached the session record, the inbox and the re-prompt messages verbatim. go test -trimpath ./internal/adapter/... and ./internal/api -run 'SessionRun|UsageLearn|SessionInterview|SessionAsk' pass.

## Links

### Out

- verifies [adapter-reads-api-error](../capability/rhiza/workflows/workflow-runs/run-resume/adapter-reads-api-error.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
