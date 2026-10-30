# review-fixers-in-rex-verified

evidence · `01M66CCEFMN3T519KE4N0R22QR` · actor [claude-code-5a92b8ed](../actor/claude-code-5a92b8ed.md)

## result

PASS

## statement

At b9affc3, ran go test ./internal/api -run TestProgramStatus -count=1 -v: TestProgramStatusReviewFixers passes. It runs a review whose check alpha confirms one leaf finding and one change-set finding through the fake Rex runner; the report trace shows, in order, leaf/1 'leaf-a: fixer working (attempt 2)', leaf/1 cleared, leaf/1 'leaf-a: verifier working (attempt 2)', fix 'review: fixer working (attempt 2)', then leaf/1 and fix cleared, with root records 'step review ...' before each. Read the diff: fixStatus in internal/run/status.go adds unjudged leaf fixer/verifier and change-set fixer sessions to Live, rexLeafID maps the change-set fixer to 'fix', and endedLeaves now matches on sameLive so leaf/n and check/n records do not mask each other's clears. make check exits 0.

## Links

### Out

- verifies [review-fixers-in-rex](../capability/rhiza/workflows/workflow-runs/run-program-status/review-fixers-in-rex.md)

## History

- `01M602TGXM184GVWJG4NE98P6Y` Rex program status shows a review run's live fixer, verifier and change-set fixer sessions, its open questions and its code check, as rhiza run status already does
