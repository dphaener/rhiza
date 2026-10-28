# review-workflow-file-verified

evidence · `01M602MA23YTZ4RH46ZWBYD11V` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

At fd593c2: review-workflow.txtar runs `rhiza run start --workflow review` on a change set built on main (builder commits with Leaf: trailers plus verifier Evidence per leaf): the code step passes; cycle 1 runs compliance, quality, qa, scope and model on opus; quality reports a major on alpha and its critic confirms it (verdict confirmed); cycle 2 runs alpha's fixer then alpha's verifier, then only quality re-runs and passes (want.log cmp plus per-check cycles from run status --json: every other check has cycle 1 only). The fake qa check writes alpha-qa-1 Evidence; the engine commits it alone as 'chore(rhiza): record qa check writes' with an Actor trailer equal to the qa session's own RHIZA_ACTOR and no Leaf: trailer, and the fixer's commit holds only alpha.txt. rhiza run build-check exits 0, rhiza run approve ends the run succeeded at signoff, and rhiza change show reports state accepted, not landed; fold-ok. TestScripts/describe-golden passes at fd593c2, so the golden is current (no request/response type changed). builtin/review.toml: code command step with no on.fail, a review checks step with the five checks and critic/fixer/verifier all on opus, and signoff on.pass=done, on.changes=review, with no land step. A built binary dry-runs `run start --workflow review` in a scratch repo, current step code. TestReviewBuiltin, TestCommitCheckWrites/None, TestChecksStepCommitsWrites and the full suite pass.

## tests

Commit fd593c2, in a detached worktree, each run with env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT: go test ./cmd/rhiza -run 'TestScripts/review-workflow' -v -count=2; go test ./internal/workflow -run TestReviewBuiltin -v -count=1; go test ./internal/api -run 'TestCommitCheckWrites' -v -count=1; go test ./internal/run -run TestChecksStepCommitsWrites -v -count=1; go test ./cmd/rhiza -run 'TestScripts/(describe-golden|review-|build-)' -count=1; go test ./... -count=1; go build -o $TMPDIR/rhiza-vw ./cmd/rhiza then rhiza-vw run start --workflow review --permission-mode manual --dry-run in a scratch repo. Run by claude-code-0e5a9d92-verifier on 2026-10-27.

## Links

### Out

- verifies [review-workflow-file](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-workflow-file.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
