# done-captures-answers-verified

evidence · `01M658PK46CX3E402TNKG7EWZY` · actor [claude-code-eb70dff3](../actor/claude-code-eb70dff3.md)

## result

PASS

## statement

Verified at 2257c5f. make check passes (exit 0). go test ./cmd/rhiza -run 'TestScripts/(session-done-captures|session-run-tokens|describe-golden)$' -count=1 passes: session-done-captures.txtar has a transcript with picked and typed answers; done lists and proposes them (picked Blue with no quote, typed DuckDB with the owner's quote); a dry run proposes nothing; a second done says already with no decisions and the change set still holds 4 node.create ops; a later capture skips all four; a session with no transcript says done with the capture warning and keeps its mark. session-run-tokens now expects the capture warning before the token-use warning, which fixes attempt 1's finding f1. go test ./internal/api -run TestSessionDone -count=1 passes, including TestSessionDoneCapturesBeforeMark: the beforeDoneMark hook sees the mark unset at capture, and the mark is kept after a failed capture. The script stands in for the exit capture with session interview, so I also ran the real path by hand with a scratch build, a fake rex and a fake claude: rhiza run start ran the fake agent in the worktree, it wrote the transcript and ran rhiza session done, which proposed answer-*-1..4 with no warnings. session run's exit capture then proposed nothing new, the change set held 4 node.create ops and a follow-up session interview skipped all four as already proposed. In the diff, session_done.go calls captureAfterRun before st.Update writes the mark, returns before capturing on a dry run, and adds capture failures to Warnings without failing.

## Links

### Out

- verifies [done-captures-answers](../capability/rhiza/lineage/interview-capture/done-captures-answers.md)

## History

- `01M61GW89J8Y62X62RNJ8Z8HAD` A plan run captures the interview's picks and answers as Decisions when the interview step says done, not only when its session exits.
