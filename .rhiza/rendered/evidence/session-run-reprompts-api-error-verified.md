# session-run-reprompts-api-error-verified

evidence · `01M7AC5CBJ0EKSSQ31J71SWTHC` · actor [claude-code-95f07963](../actor/claude-code-95f07963.md)

## result

PASS

## statement

At 8f6535fa8: go test -trimpath ./internal/api/ -run SessionRunAPIError passed (7 tests, also -count=3); the Checks: block re-ran with exit 0 (go test ./internal/api -run 'SessionRunAPIError|SessionRun|UsageLearn|SessionInterview|SessionAsk', go test ./internal/adapter/..., rhiza ship generate, rhiza render); gofmt and go vet clean, no TODOs. Diff read: watchLimit now polls apiErrorWatch on the transcript past base; autonomous sessions get a session.message naming the error after the backoff, with a plain idle hook report first (asserted at the first send; later sends find it already idle), cancelled by a later assistant line, at most 3 then recorded with the 'stopped on an API error: ' detail; attended sessions recorded at once with no message; usage-limit pause unchanged; errors before base ignored. Inbox checks use runAPIErrorAll, the func inboxItems is given. Only the three named files plus the builder's actor record changed.

## Links

### Out

- verifies [session-run-reprompts-api-error](../capability/rhiza/workflows/workflow-runs/run-resume/session-run-reprompts-api-error.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
