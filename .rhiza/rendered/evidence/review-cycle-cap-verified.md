# review-cycle-cap-verified

evidence · `01M5ZZ1AESDAB5A7P04X6CHCZA` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

Re-ran the leaf's tests at 0245da3 in a detached worktree; all pass. Per clause: (1) check failing three times -> exactly two fix cycles: TestCycleCapRaisesQuestions asserts leaf sessions are only fixer/verifier of cycles 2 and 3 and len(Cycles)==3; (2) question blocking the named leaf: same test asserts one AddLeafQuestion (no review question) naming leaf-a/alpha/f1 and a review-question entry Leaf=L1, Question=QL1; api TestAddLeafQuestion asserts that AddLeafQuestion writes a blocks link from the question to the leaf; (3) no sign-off step entered: no 'signoff' command and no entered entry for step signoff; (4) run failed: Result==Failed, outcomes [review:fail], detail names the unanswered question; (5) after an answering Decision, resume runs one more cycle with the Decision in the fixer's prompt: TestCycleCapAnsweredResume (unanswered resume runs nothing; once answered, cycle 4 is an answered cycle with Decision rename-it, the cycle-4 fixer LeafSpec.Answers carries the Decision and statement, sign-off then runs) plus api TestCycleCapFixerPrompts asserting the rendered leaf and change-set fixer prompts contain '- Decision rename-it answers question ...: <statement>' and a within-cap fixer's does not; (6) autonomous: TestCycleCapAutonomous asserts arbiters==[arbiter QL1], the answered cycle 4 runs in the same visit with the arbiter's Decision in the fixer's answers, the cycle-4 failure ends the run Failed with sign-off not run, and a retry neither arbitrates again nor runs another cycle. Also passed: TestCycleCapQuestionsWithoutLeaf, TestAddReviewQuestion, TestCycleCapStatus, TestChecksStepCycleLimit, TestFixLeafNoCommitFailsCycle, TestFixChangeSetNoCommitFailsCycle, the review-/build-/run- testscripts and go test ./... (all packages ok).

## tests

Commit 0245da3, detached worktree $TMPDIR/rv-review-cycle-cap, each command prefixed env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT: go test ./internal/run -run 'TestCycleCap|TestChecksStepCycleLimit|TestFixLeafNoCommitFailsCycle|TestFixChangeSetNoCommitFailsCycle' -v -count=1; go test ./internal/api -run 'TestAddReviewQuestion|TestCycleCap' -v -count=1; go test ./cmd/rhiza -run 'TestScripts/(review-|build-|run-)' -count=1; go test ./... -count=1. Run by claude-code-0e5a9d92-verifier on 2026-10-27.

## Links

### Out

- verifies [review-cycle-cap](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-cycle-cap.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
