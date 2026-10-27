# plan-gate-verified

evidence · `01M5XMT0CM0BE90MT9XVR93T2G` · actor [claude-code-c3d71038-verifier](../actor/claude-code-c3d71038-verifier.md)

## result

PASS on 2026-10-26 at commit 55568ab on branch cs/the-plan-workflow, re-run with -count=1 by an independent verifier (claude-code-c3d71038-verifier) after b8bdc34

## statement

At HEAD, after b8bdc34: testscript run-gate passes: a normal run whose landed record adds an intent served by two drafts proposes both, accepts none, leaves the record landed, and run status --json names rhiza accept --rule owner for each; autonomous with a clean review both are accepted with rule autonomous and a reason naming the run; autonomous with a round-3 question both stay proposed and the gate exits 1 naming the question, then after a Decision answers it both are accepted; a run on an intention change set serving an existing intent proposes that change set; no RHIZA_RUN is usage. plan-workflow-autonomous-left: with review questions left in the draft and in a landed record write, each arbiter answers in the question's home change set (the draft) or its own landed record change set, and the gate then accepts the draft with rule autonomous. autonomous-gate-rule and plan-workflow-autonomous pass. Hands-on: a scratch build of HEAD, all RHIZA_* unset, in a scratch repo: rhiza run gate exited 2 saying RHIZA_RUN is not set.

## tests

go test -count=1 ./internal/run ./internal/api; go test -count=1 ./cmd/rhiza -run 'TestScripts/(run-arbitrate-asked|arbiter-step|arbiter-question|arbiter-answer|run-gate|autonomous-gate-rule|plan-workflow-autonomous|plan-workflow-autonomous-left)$'; hands-on: rhiza run gate with a scratch build in a scratch repo with no RHIZA_RUN (exit 2)

## Links

### Out

- verifies [plan-gate](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
