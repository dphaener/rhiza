# session-run-reprompts-api-error-verified-2

evidence · `01M7AHPRAAY7602HE75ARAM5N5` · actor [claude-code-d513d9e3](../actor/claude-code-d513d9e3.md)

## result

PASS

## statement

Verified attempt 2 at 706598253. Re-ran: go vet ./internal/api/ (exit 0); go test -trimpath ./internal/api/ -run 'TestAPIErrorWatch|TestSessionRunAPIError|TestSessionRun' -count=2 (exit 0); acceptance go test -trimpath ./internal/api/ -run SessionRunAPIError -v: all 7 TestSessionRunAPIError* pass (reprompt after backoff with plain idle record and no inbox item, cancel by assistant line, gives up after three consecutive errors leaving the APIErrorDetail record, re-prompts after recovery, attended records at once with no message, error before base is nothing, usage limit keeps pause with no message); rhiza ship generate and rhiza render exit 0 and leave the tree clean; make self-check exit 0. Read the diff: poll now starts each wait when a new error past the send's transcript length is seen, a fourth error goes to the inbox at once, recordIdle reports write failures and stop retries on the next poll; covered by the new TestAPIErrorWatchReportsFailedRecord/FourthErrorAtOnce/WaitsFromErrorSeen. The commit touches only session_run_apierror.go, its test and .rhiza records; no TODOs or stubs; fakes only. The full go test ./internal/api/... could not finish here: on a host at load 11-13 it hit the test timeout twice with no failing test, running a different unrelated test each time (TestRootIdlePoll, its production-sized fixture, then TestDaemonAutoUpdateRelease), so it is left to the PR's CI.

## Links

### Out

- verifies [session-run-reprompts-api-error](../capability/rhiza/workflows/workflow-runs/run-resume/session-run-reprompts-api-error.md)

## History

- `01M7A949Q4VGB6D4PPY9A5GZEK` Fix Run sessions that hit an API error stop at the prompt and sit idle forever, unflagged. On 2026-11-13 around 09:47-09:52Z three agent sessions ended their turn with Claude Code's 'API Error: 403 Attention Required! | Cloudflare' (transient): verifier 4594e9fa (run 01M79TRE3W, leaf daemon-holds-steps-on-stale-guide, idle 09:47:45-10:05) and review checks quality dbb3e57f and compliance 8b78600b (run 01M7A144DD, cycle 1, idle 09:51-10:05). Their processes stayed alive, the runs showed 'running', no inbox item, no watchdog action; load fell to ~2 while three slots did nothing. Unstuck by the orchestrator: rhiza session message <session> asking each to carry on; all resumed within seconds. Wanted: the run (or watchdog) notices a session whose transcript ends in an 'API Error:' assistant message (or whose transcript has not grown for N minutes with no child process) and re-prompts it automatically, retrying a few times with backoff, then raises an inbox item. Related, different cause: 01M6N9AB3X (idle after background shell), 01M79MDD5R (idle after unknown slash command).
