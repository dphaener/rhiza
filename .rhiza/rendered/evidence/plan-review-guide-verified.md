# plan-review-guide-verified

evidence · `01M5XMSZJF7E1F13XK6NX7BQYA` · actor [claude-code-c3d71038-verifier](../actor/claude-code-c3d71038-verifier.md)

## result

PASS on 2026-10-26 at commit 5fff4c7 (code as of 2963c0f) on branch cs/the-plan-workflow, re-run with -count=1 by an independent verifier (claude-code-c3d71038-verifier)

## statement

TestPlanReviewText finds in the plan-review topic the five kinds (missing-acceptance, missing-provenance, unreal-depends-on, outside-intent, contradiction), rhiza run findings and rhiza run verdict with their JSON file formats, how the plan's change sets are found ('serving an intent it adds', rhiza change show, rhiza slice, rhiza why) and 'Never add a finding'; TestPlanKinds finds it a topic; the cap test keeps it within TopicCap; testscript guide shows guide ls listing it. Hands-on: the scratch build's guide ls listed plan-review as a topic of 499 tokens and guide plan-review printed all five kinds and 'Never add a finding'.

## tests

go test -count=1 ./internal/guide -run 'TestPlanReviewText|TestPlanKinds|Cap'; go test -count=1 ./cmd/rhiza -run 'TestScripts/guide$'; hands-on: rhiza guide ls; rhiza guide plan-review | grep

## Links

### Out

- verifies [plan-review-guide](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-guide.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
