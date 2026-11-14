# leafless-refusal-by-leaves-step-verified

evidence · `01M7C167VJSAP9023DP2G5V83A` · actor [claude-code-77fa7410](../actor/claude-code-77fa7410.md)

## result

PASS

## statement

At 417450213: go test -trimpath ./internal/api -run TestCheckShipLeaves exit 0 (leafless change set refused with exit 3 naming the workflow for ship-change-classic, a repo leaves-step workflow and a repo workflow including build; accepted for a workflow with no leaves step); go test -trimpath ./cmd/rhiza -run TestScripts/run-ship-needs-leaves exit 0 (script pinned to ship-change-classic, a repo ship-change without leaves step and the conductor start on a leafless change set); rhiza ship generate and rhiza render exit 0 with no changes; make self-check exit 0 in a clean detached worktree at HEAD (in the run checkout it fails only on the verifier's own uncommitted actor record). Diff read: run.go gates checkShipLeaves on workflowHasLeavesStep(def) over expanded steps instead of the ship-change name; only the three named files plus the builder's actor record and render changed; no stubs or TODOs.

## Links

### Out

- verifies [leafless-refusal-by-leaves-step](../capability/rhiza/workflows/built-in-workflows/leafless-refusal-by-leaves-step.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
