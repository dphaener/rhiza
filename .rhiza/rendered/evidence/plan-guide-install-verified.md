# plan-guide-install-verified

evidence · `01M5XMSXPTW024R9B5QGPJSMRT` · actor [claude-code-c3d71038-verifier](../actor/claude-code-c3d71038-verifier.md)

## result

PASS on 2026-10-26 at commit 5fff4c7 (code as of 2963c0f) on branch cs/the-plan-workflow, re-run with -count=1 by an independent verifier (claude-code-c3d71038-verifier)

## statement

Testscript guide-install, with its own HOME, finds skills/rhiza-plan/SKILL.md after guide install, its frontmatter naming rhiza-plan with the plan text's description, and its body (after the frontmatter) byte-equal to rhiza guide plan; after the file is removed status reports the install out of date with next rhiza guide install, and a second install restores it as current. Testscript status repeats this for --json (current false) and --hook; TestStatusGuideFiles checks every installed file, missing or changed, makes the install not current. Hands-on: in a scratch HOME the scratch build's guide install wrote skills/rhiza-plan/SKILL.md with name: rhiza-plan; after deleting it status said '(out of date)' with next rhiza guide install, and a second install wrote only that file and status said '(current)'.

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/guide-install$'; go test -count=1 ./cmd/rhiza -run 'TestScripts/status$'; go test -count=1 ./internal/api -run TestStatusGuideFiles; hands-on: guide install, rm the rhiza-plan SKILL.md, status, guide install, status in a scratch HOME

## Links

### Out

- verifies [plan-guide-install](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-guide-install.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
