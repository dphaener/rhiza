# plan-workflow-file-verified

evidence · `01M5XMT0NATVRAQAYYPYXME13C` · actor [claude-code-c3d71038-verifier](../actor/claude-code-c3d71038-verifier.md)

## result

PASS on 2026-10-26 at commit 5fff4c7 (code as of 2963c0f) on branch cs/the-plan-workflow, re-run with -count=1 by an independent verifier (claude-code-c3d71038-verifier)

## statement

TestPlanBuiltin finds plan in BuiltinNames, parses its file into interview, review, critic, check, revise, gate; the interview prompt names rhiza-plan and rhiza guide plan; the review, critic and revise prompts use only the {change} placeholder and name rhiza guide plan-review and never a plan file; check runs rhiza run review-check (fail -> revise) and gate runs rhiza run gate. Testscripts plan-workflow and plan-workflow-autonomous run the built-in with fake agents to the gate, and git status in the worktree shows only .rhiza/ paths. Hands-on: read internal/workflow/builtin/plan.toml (no plan file named) and the scratch build's run start --workflow plan --permission-mode manual --dry-run would start at step interview.

## tests

go test -count=1 ./internal/workflow -run TestPlanBuiltin; go test -count=1 ./cmd/rhiza -run 'TestScripts/plan-workflow$'; go test -count=1 ./cmd/rhiza -run 'TestScripts/plan-workflow-autonomous$'; hands-on: rhiza run start --workflow plan --permission-mode manual --dry-run

## Links

### Out

- verifies [plan-workflow-file](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-file.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
