# plan-decisions-from-source-verified

evidence · `01M5XMSZ1BG4TS2A51APK78SQE` · actor [claude-code-c3d71038-verifier](../actor/claude-code-c3d71038-verifier.md)

## result

PASS on 2026-10-26 at commit 5fff4c7 (code as of 2963c0f) on branch cs/the-plan-workflow, re-run with -count=1 by an independent verifier (claude-code-c3d71038-verifier)

## statement

Testscript session-interview-source adds a Source with --raw-session for a fake session whose transcript holds one picked and one typed answer, then runs session interview: answer-1 has picked Blue and wording, answer-2 the quote 'Something like DuckDB, but local', both link derived-from the Source (two link.create ops), a second run adds nothing; with no such Source both are proposed with empty links_out. session-interview-landed covers the same after the record landed (a landed record change set). TestPlanText (plan-decisions-from-source) finds derived-from and 'quotes the owner's typed words'. Hands-on: the scratch build's guide plan carried the derived-from rule (counted among 12 matching lines).

## tests

go test -count=1 ./cmd/rhiza -run 'TestScripts/session-interview-source$'; go test -count=1 ./cmd/rhiza -run 'TestScripts/session-interview-landed$'; go test -count=1 ./internal/guide -run 'TestPlanText'; hands-on: rhiza guide plan | grep derived-from

## Links

### Out

- verifies [plan-decisions-from-source](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-decisions-from-source.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
