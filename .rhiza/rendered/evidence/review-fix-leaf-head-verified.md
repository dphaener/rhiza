# review-fix-leaf-head-verified

evidence · `01M6076KSZW0968TMYM56G0XFM` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

At d60c3c3, after fix 790f2f1 (leaf fixer prompt now reads 'Confirmed findings to fix in review cycle N on leaf <leaf>:'), TestFixLeafEngine and TestFixLeafEngineNoCommit pass against a real git repo: the fixer's prompt names leaf-a and the finding, its new commit carries Leaf: leaf-a, leaf-a's verifier starts on that commit and its Evidence commit carries the verifier's actor, only the failed check alpha re-runs; a fixer with no commit naming the leaf fails the cycle and starts nothing further. TestOwnerCycleFixerPrompts, TestCycleCapFixerPrompts and the internal/run TestFixLeaf*/TestFixActions/TestNextCheckActionsFixPhase tests pass; go test ./... passes.

## tests

env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./internal/api -run 'TestFixLeafEngine|TestOwnerCycleFixerPrompts|TestCycleCapFixerPrompts' -v -count=1; env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./internal/run -run 'TestFixLeaf|TestFixActions|TestNextCheckActionsFixPhase' -v -count=1; env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT go test ./... -count=1; commit d60c3c3 (detached worktree); run by claude-code-0e5a9d92-verifier on 2026-10-28

## Links

### Out

- verifies [review-fix-leaf](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-fix-leaf.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
