# question-session-panes-close-verified

evidence · `01M68J1FDQYSDKDZ0X6Q4GD0HF` · actor [claude-code-a62c22b9](../actor/claude-code-a62c22b9.md)

## result

PASS

## statement

At e086794 ran go test ./internal/api -run 'TestArbiterPaneClosesAfterVerdict|TestRunRecommendClosesAdvisorPane' -v -count=1: all pass. TestArbiterPaneClosesAfterVerdict drives an autonomous run through run.Engine with the stepRunner's own ClosePane and a fake rex: a passing verdict gives exactly one `rex block close` of block:arb, and the verdict is already journalled when the close happens. A failed verdict and an arbiter lost while still running give no close, and no close happens while the arbiter is being polled. TestRunRecommendClosesAdvisorPane: run recommend journals the recommendation, then closes block:a2 once, and a dry run closes nothing. go test ./internal/run -run TestArbitrate passes. Mutation in a scratch copy: removing r.closeQuestionPane(session) in engine.go and closeSessionPane in run_recommend.go makes both tests fail. make fmt-check vet test pass. make check's self-check flagged only the unrendered actor file of this verifier session, which is not part of the leaf commit. This attempt fixes the attempt-1 findings (fake-rex arbiter test, still-running case, ordering, stale doc comment).

## Links

### Out

- verifies [question-session-panes-close](../capability/rhiza/runtime/rex-sessions/question-session-panes-close.md)

## History

- `01M66M3S5PZD9PXDP0X3N1ZTXH` A run's agent-step, arbiter and advisor sessions open as tiled panes of the run's window, like its leaf, check, critic and fixer sessions, so autonomous runs with many arbiter calls add no window each
