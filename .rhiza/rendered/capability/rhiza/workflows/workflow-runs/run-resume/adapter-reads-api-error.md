# adapter-reads-api-error

capability · `01M7A9C4BPN1BANYKJJREMJNFR` · parent [run-resume](../run-resume.md)

## acceptance

A Go test in internal/adapter/claudecode (go test -trimpath ./internal/adapter/...): a transcript whose last assistant line is the synthetic 'Please run /login · API Error: 403 Attention Required! | Cloudflare' message reports an API error with that text; one whose last assistant line is the usage-limit message reports none (and UsageLimit still reports the limit); an API-error line followed by a later model assistant line reports none; 'API Error' in the model's own prose or a tool call's input reports none; an empty transcript reports none. The api package's fake adapters gain the method and go test -trimpath ./internal/api/ -run 'SessionRun|UsageLearn|SessionInterview|SessionAsk' passes. Touches internal/adapter/adapter.go, internal/adapter/claudecode/claudecode.go, internal/adapter/claudecode/claudecode_test.go, internal/api/*_test.go (fake adapters)

## statement

The agent adapter tells from a transcript that the agent's last turn ended on an API error that is not its usage limit, and returns the error's text: for Claude Code, the last assistant line is the synthetic API-error message (isApiErrorMessage or model <synthetic>) whose text is not the usage-limit message

## Links

### Out

- decided-by [api-error-reprompt-in-session-run](../../../../../decision/api-error-reprompt-in-session-run.md)

### In

- depends-on [session-run-reprompts-api-error](session-run-reprompts-api-error.md)
- verifies [adapter-reads-api-error-verified](../../../../../evidence/adapter-reads-api-error-verified.md)
- verifies [adapter-reads-api-error-qa-1](../../../../../evidence/adapter-reads-api-error-qa-1.md)
- verifies [pr-548-fix-run-sessions-that-hit-an](../../../../../evidence/pr-548-fix-run-sessions-that-hit-an.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
