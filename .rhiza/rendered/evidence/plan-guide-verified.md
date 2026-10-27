# plan-guide-verified

evidence · `01M5XMSXEA9NQVB2QJYRTRVYWA` · actor [claude-code-c3d71038-verifier](../actor/claude-code-c3d71038-verifier.md)

## result

PASS on 2026-10-26 at commit 5fff4c7 (code as of 2963c0f) on branch cs/the-plan-workflow, re-run with -count=1 by an independent verifier (claude-code-c3d71038-verifier)

## statement

TestPlanKinds finds plan as a skill (# rhiza-plan, ## Steps) and plan-review as a topic; TestPlanText finds in the plan text rhiza apply, a statement= and acceptance= recipe with a decided-by link, 'no plan file' and 'no checkboxes'; TestCaps and TestDescriptions keep every text within SkillCap/TopicCap and every description within MaxDescription; testscript guide shows guide ls listing '- plan skill ... Plan with the owner' and guide plan printing the decided-by apply line. Hands-on: the scratch build's guide ls listed plan as a skill (1517 tokens, under SkillCap 2000) and guide plan held the decided-by recipe and 'no plan file'.

## tests

go test -count=1 ./internal/guide -run 'TestPlan|TestDescriptions|Cap'; go test -count=1 ./cmd/rhiza -run 'TestScripts/guide$'; hands-on: rhiza guide ls; rhiza guide plan | grep

## Links

### Out

- verifies [plan-guide](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-guide.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
