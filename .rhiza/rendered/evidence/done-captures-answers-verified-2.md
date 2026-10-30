# done-captures-answers-verified-2

evidence · `01M659J54S2XGMDR9VT4D69M6Z` · actor [claude-code-583d38c0](../actor/claude-code-583d38c0.md)

## result

PASS

## statement

At 4a4bf55 (attempt 2, fixing review finding f1), I ran go test ./cmd/rhiza -run 'TestScripts/(session-done-captures|session-run-tokens)$' -count=1 and go test ./internal/api -run Done -count=1, and both passed. Then make fmt-check vet test had no FAIL. In session-done-captures.txtar a fake claude now runs under rhiza session run. It writes a transcript with picked and typed answers, and a dry-run done captures nothing (show exits 4, 'no node or link matches'). session done --json lists 4 answer Decisions, which show reads back while the agent still runs (picked Blue; typed quote 'Something like DuckDB, but local'). A second done returns already:true with no decisions. The exit capture in session run (captureAfterRun; the CLI sets Home, so it runs) emits no Warning, and the change set still has exactly 4 node.create. A no-transcript session says done with a warning and keeps its mark. TestSessionDoneCapturesBeforeMark uses the beforeDoneMark hook to show the mark is not set at capture. make self-check fails only on stale .rhiza/rendered from the run's chore(rhiza) record commits. This leaf commit does not touch .rhiza.

## Links

### Out

- verifies [done-captures-answers](../capability/rhiza/lineage/interview-capture/done-captures-answers.md)

## History

- `01M61GW89J8Y62X62RNJ8Z8HAD` A plan run captures the interview's picks and answers as Decisions when the interview step says done, not only when its session exits.
