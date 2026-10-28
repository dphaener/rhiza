# review-fix-leaf-verified

evidence · `01M5ZX23KS2GYH9ZWXA9N8YJG8` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

Re-ran the leaf's tests at d145f9a in a detached worktree with RHIZA_* unset; all pass. internal/api/run_fixleaf_test.go TestFixLeafEngine runs the engine over a real git repo (evidenceCommitRunner: git init + rhiza init, real commit/trailer judging) and asserts each clause: for confirmed major f1 on leaf-a the fixer prompt contains 'fix leaf leaf-a in cycle 2', 'alpha f1', 'major', 'fix.go:1', the reason and '`Leaf: leaf-a`'; the fixer's new commit has Leaf trailer leaf-a and its judgement passes on that commit; the verifier spec Commit and its checkout HEAD equal the fix commit (checkout removed after); the verifier judgement is 'Evidence written by the verifier' and the Evidence commit's Actor trailer equals the verifier actor's name; check sessions are exactly alpha check 1, beta check 1, alpha critic 1, alpha check 2 (only the failed check re-runs) and the verifier is judged before cycle 2's check starts; no leaf outcome entry. TestFixLeafEngineNoCommit asserts a fixer with no commit naming leaf-a fails: step fails, no verifier starts, no cycle-2 check session, fixer judgement 'no new commit since ...', alpha cycle 2 failed with 'cycle 2: fix failed: leaf leaf-a: fixer failed: ...'. Unit tests TestFixLeafOrder/NoCommitFailsCycle/VerifierFailsCycle/NoneWithoutLeafFindings, TestFixActions, TestNextCheckActionsFixPhase pass. git show d145f9a (fixActions/fixFindings in internal/run/checks.go, fixerPrompt in internal/api/run_leaves.go) matches the acceptance.

## tests

env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./internal/api -run 'TestFixLeafEngine' -v -count=1; same env go test ./internal/run -run 'TestFixLeaf|TestFixActions|TestNextCheckActionsFixPhase' -v -count=1; same env go test ./cmd/rhiza -run 'TestScripts/(review-|build-)' -count=1; same env go test ./... -count=1 (all ok); commit d145f9a; run by claude-code-0e5a9d92-verifier on 2026-10-27

## Links

### Out

- verifies [review-fix-leaf](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-fix-leaf.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
