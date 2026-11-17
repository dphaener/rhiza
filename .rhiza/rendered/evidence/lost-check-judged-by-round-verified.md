# lost-check-judged-by-round-verified

evidence · `01M7MH41JDVR5FMKY5AX1A4H49` · actor [claude-code-5013a51b](../actor/claude-code-5013a51b.md)

## result

PASS

## statement

At ed7678fe4 re-ran the Checks: block: go test -trimpath ./internal/run -run 'TestChecksStep|TestCritic|TestNextCheckActions' -count=1 (ok), go test -trimpath ./internal/run/... (ok), gofmt -l internal/run (empty), go vet ./internal/run/ (ok), rhiza ship generate (exit 0). The five acceptance cases exist in internal/run/checks_test.go and pass (-v): TestChecksStepLostWithRoundsPass (all lost with empty rounds: each judged pass, no cycle 2 entry, no session started), TestCriticLostWithVerdictsPass (major finding starts alpha's critic; lost critic with verdict judged pass, check passed by 'critic dropped f1'), TestChecksStepLostNoRoundFails (no round: fail, cycle 2 runs alpha alone), TestCriticLostMissingVerdictFails (missing verdict on blocker: critic fail, 'critic failed'), TestChecksStepKnownFailWithRoundFails (non-lost exit 1 with round: fail). Existing TestChecksStep*/TestCritic* tests unchanged and pass. Diff touches only checks.go and checks_test.go (plus the actor record/render); lostWorkRecorded applies only when status.Done && Lost && !Environment and only for check/critic roles, requires the round (check) or a verdict on every serious finding (critic); AgentStatus, internal/api, cycle cap and later-cycle check selection are untouched. No stubs or TODOs.

## Links

### Out

- verifies [lost-check-judged-by-round](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-checks-step/lost-check-judged-by-round.md)

## History

- `01M7MGF8J9MA30XXQ8GTR4KXRW` Fix Resuming a classic ship-change run whose check sessions died in a crash marks the whole review cycle failed and opens a new cycle of all four checks, even when every check's round was later recorded and rhiza run review-check passes
