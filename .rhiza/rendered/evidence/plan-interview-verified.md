# plan-interview-verified

evidence · `01M5XMSY7W1NR7B5EHP4CPPM3F` · actor [claude-code-c3d71038-verifier](../actor/claude-code-c3d71038-verifier.md)

## result

PASS on 2026-10-26 at commit 5fff4c7 (code as of 2963c0f) on branch cs/the-plan-workflow, re-run with -count=1 by an independent verifier (claude-code-c3d71038-verifier)

## statement

TestPlanText (plan-interview) finds rhiza slice and rhiza why, 'Never ask a fact', 'One topic per message', recommendation, AskUserQuestion, 'Only typed words are quoted', the wording and quote fields, and the question recipe with a blocks link; the cap test keeps it within SkillCap. Hands-on: the scratch build's guide plan printed 'Never ask a fact', 'Only typed words are quoted', rhiza slice and AskUserQuestion.

## tests

go test -count=1 ./internal/guide -run 'TestPlanText|Cap'; hands-on: rhiza guide plan | grep

## Links

### Out

- verifies [plan-interview](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-interview.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
