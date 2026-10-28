# review-check-worktree-verified

evidence · `01M5ZVA8QVRZXCHCDWXGEDNFZQ` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

Re-ran at 583b957 in a fresh detached worktree with RHIZA_* unset; all pass. review-check-worktree.txtar asserts each clause: (1) detached worktree at branch head: fake claude logs dir=checkout head=branch-head detached=yes root=run for alpha-1, beta-1, beta-2 (cmp want.log) and each pwd is under .git/rhiza/checkouts/<run>/<check>-check-<cycle>; (2) gone after session whatever its outcome: alpha-1 and beta-2 pass, beta-1 exits 1, and none of the three pwds exist afterwards nor appear in git worktree list; (3) finding and Evidence in the run worktree: alpha's findings entry (check alpha, cycle 1) is in the run journal and Evidence alpha-ok is a .json under the run worktree's .rhiza/log, while the checkout's git status stayed clean; (4) staged, unstaged and untracked edits made by a prior command step in the run worktree compare byte-equal before/after (git diff, diff --cached, status). Engine tests TestChecksStepCheckouts, TestChecksStepCheckoutResume, TestChecksStepCheckoutErrors, TestChecksStepNoFindingsRound pass (open/remove on pass and fail, resume, start/remove failures, D10 no-findings-round). Full go test ./... passes (21 packages ok).

## tests

env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./cmd/rhiza -run 'TestScripts/review-check-worktree' -v -count=1; same prefix go test ./internal/run -run 'TestChecksStepCheckout|TestChecksStepNoFindingsRound' -v -count=1; same prefix go test ./... -count=1; commit 583b957; run by claude-code-0e5a9d92-verifier on 2026-10-27

## Links

### Out

- verifies [review-check-worktree](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-check-worktree.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
