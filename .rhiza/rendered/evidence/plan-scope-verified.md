# plan-scope-verified

evidence · `01M5XMSXZAF1PJZW7SKJ0P4H31` · actor [claude-code-c3d71038-verifier](../actor/claude-code-c3d71038-verifier.md)

## result

PASS on 2026-10-26 at commit 5fff4c7 (code as of 2963c0f) on branch cs/the-plan-workflow, re-run with -count=1 by an independent verifier (claude-code-c3d71038-verifier)

## statement

TestPlanText (plan-scope) finds in rhiza guide plan both scopes (`intention` and `record`), milestone, the serves rule, child intents, rhiza change new, 'intent only' for change sets not planned now, rhiza land of the record, 'before review', depends-on and drafts left for the gate; the cap test keeps the text within SkillCap. Hands-on: the scratch build's guide plan printed 'child intents' and rhiza change new, at 1517 tokens.

## tests

go test -count=1 ./internal/guide -run 'TestPlanText|Cap'; hands-on: rhiza guide plan | grep

## Links

### Out

- verifies [plan-scope](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-scope.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
