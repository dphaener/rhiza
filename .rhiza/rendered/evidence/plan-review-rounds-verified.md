# plan-review-rounds-verified

evidence · `01M5XMSZV2WCFFNCNCKYD8WN8F` · actor [claude-code-c3d71038-verifier](../actor/claude-code-c3d71038-verifier.md)

## result

PASS on 2026-10-26 at commit 5fff4c7 (code as of 2963c0f) on branch cs/the-plan-workflow, re-run with -count=1 by an independent verifier (claude-code-c3d71038-verifier)

## statement

Testscript run-review-check in a held run: no round fails; a finding confirmed in round 1 fails (exit 1), names f1 and adds no question and no asked entry; a clean round 2 passes; round 3 with two confirmed findings passes and adds review-<run>-r3-f1 (statement holding the finding and the critic's reason, blocks cap-a) and r3-f2 (on a decision, no blocks link), with two asked entries; a second check adds nothing; no RHIZA_RUN is usage. run-review-check-landed covers the questions after the record landed early. Hands-on: in a scratch autonomous held run with one finding confirmed each round, the scratch build's review-check exited 1 in round 1 and in round 2 (the round-2 failure the testscript does not show) and 0 in round 3, adding question review-...-r3-f1 and one asked entry.

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/run-review-check$'; go test -count=1 ./cmd/rhiza -run 'TestScripts/run-review-check-landed$'; hands-on: rhiza run findings/verdict/review-check over three rounds in a scratch held run

## Links

### Out

- verifies [plan-review-rounds](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-rounds.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
