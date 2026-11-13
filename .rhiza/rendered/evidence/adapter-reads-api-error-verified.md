# adapter-reads-api-error-verified

evidence · `01M7AB4T2EW9SJ7N7R41JWDXY2` · actor [claude-code-87ac5975](../actor/claude-code-87ac5975.md)

## result

PASS

## statement

At e453aa321 re-ran the Checks block: go test -trimpath ./internal/adapter/... (ok, TestAPIError passes with -count=1), go test -trimpath ./internal/api/ -run 'SessionRun|UsageLearn|SessionInterview|SessionAsk' (ok), rhiza ship generate and rhiza render (exit 0, no tracked changes); go build ./... and go vet ./... clean. Diff adds APIError to adapter.Adapter and the Claude Code adapter (last assistant line is synthetic/isApiErrorMessage and not the usage-limit text, returns its text), sharing lastAssistant/isLimitText with UsageLimit. TestAPIError covers the Cloudflare 403 line (reports its text), the usage-limit line (none, UsageLimit still reports), an API error followed by a model line (none), 'API Error' in model prose and in tool_use input (none) and an empty transcript (none). Fake adapters in internal/adapter and internal/api tests gained the method. Only the listed files plus the builder's actor record changed; no stubs or TODOs.

## Links

### Out

- verifies [adapter-reads-api-error](../capability/rhiza/workflows/workflow-runs/run-resume/adapter-reads-api-error.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
