# plan-workflow-verified

evidence · `01M5XMSX5Q71VCM576RPGGBNSB` · actor [claude-code-c3d71038-verifier](../actor/claude-code-c3d71038-verifier.md)

## result

PASS on 2026-10-26 at commit 5fff4c7 (code as of 2963c0f) on branch cs/the-plan-workflow, re-run with -count=1 by an independent verifier (claude-code-c3d71038-verifier)

## statement

Testscript plan-workflow runs rhiza run start <record> --workflow plan --permission-mode manual on a committed draft record change set with one fake agent playing every role: the interview writes the milestone and lands the record, a fresh reviewer records one missing-acceptance finding, the critic confirms it, review-check fails (exit 1) and routes to revise, the reviser adds the acceptance, round 2 is clean, check passes and the gate runs; the run succeeds, the planned change set is proposed and holds the reviser's acceptance, csv-streaming quotes the owner's typed words and links derived-from, git status --porcelain in the worktree shows only .rhiza/ paths, and after render rhiza check prints ok. Hands-on: the scratch build's rhiza run start --workflow plan --permission-mode manual --dry-run in a scratch repo printed 'Run ... would start', workflow plan, current step interview, nothing written.

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/plan-workflow$'; hands-on: rhiza run start --workflow plan --permission-mode manual --dry-run in a scratch repo

## Links

### Out

- verifies [plan-workflow](../capability/rhiza/workflows/built-in-workflows/plan-workflow.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
