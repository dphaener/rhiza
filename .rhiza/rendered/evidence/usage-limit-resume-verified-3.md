# usage-limit-resume-verified-3

evidence · `01M65NPQHTX2MZ5XBPCD0AA7GC` · actor [claude-code-b72a0480](../actor/claude-code-b72a0480.md)

## result

PASS

## statement

At abddd6f, ran `go test ./cmd/rhiza -run TestScripts/session-run-usage-limit -count=1`: PASS. Its fake claude writes Claude Code's synthetic usage-limit error line (isApiErrorMessage, model <synthetic>, 'usage limit reached|1700000000') on its first launch and finishes on the next. The script asserts exactly one 'start --session-id <id>' and exactly one 'resume --resume <id>' for the same stored session ID, and that the session ended with exit 0. So the run resumed the stored session with the CLI's --resume flag, the step completed, and no second session started. Also passing: TestUsageLimit, TestUsageLimitIgnoresModelTurns (tool_use input and model prose holding the phrase are not the limit), TestSessionRunStopsAgentOpenAtLimit, TestSessionRunFailsWhenPausesRunOut, TestSessionHookStopAtLimitIsPaused, TestSessionRunFreshSessionWhenTranscriptGone. gofmt and vet are clean. go test ./... passes, except one timing flake in TestGeneratedModel (58ms vs 50ms under load) that passed on its own rerun. make check's self-check failed only on stale render of this verifier's own actor record. Read the diff: UsageLimit now matches only the text blocks of the synthetic API-error line.

## Links

### Out

- verifies [usage-limit-resume](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-resume.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
