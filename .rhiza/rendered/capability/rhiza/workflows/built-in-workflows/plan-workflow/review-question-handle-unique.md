# review-question-handle-unique

capability · `01M6W7QGXS2VNCJ3KMFRMW78XE` · parent [plan-workflow](../plan-workflow.md)

## acceptance

Go tests in internal/api and the run-review-check testscript(s) covering it: two runs whose IDs share their first 8 characters each reach round 3 with a confirmed finding f1 and leave questions with different handles carrying their full IDs; with the full-ID handle already taken by an unrelated node, review-check adds review-<id>-r3-f1-2, journals that handle, a second review-check adds nothing, and the autonomous gate finds the question left for the finding.

## statement

rhiza run review-check names each question it leaves review-<full run ID, lower-cased>-r<round>-<finding>, so runs started in the same second never share a handle; when that handle is already taken by a node that is not this run's question for this finding, it takes the first free one (base-2, base-3, ...), journals the handle it used in the asked entry, and the run's later checks (a repeated review-check, the autonomous accept gate) find its questions by the journalled entry instead of rebuilding the handle.

## Links

### Out

- depends-on [plan-review-rounds](plan-review-rounds.md)
- decided-by [review-question-handle-full-run-id](../../../../../decision/review-question-handle-full-run-id.md)

### In

- verifies [review-question-handle-unique-verified](../../../../../evidence/review-question-handle-unique-verified.md)
- verifies [review-question-handle-unique-qa-1](../../../../../evidence/review-question-handle-unique-qa-1.md)
- verifies [review-question-handle-unique-verified-2](../../../../../evidence/review-question-handle-unique-verified-2.md)
- verifies [pr-369-fix-review-question-handles](../../../../../evidence/pr-369-fix-review-question-handles.md)

## History

- `01M6W7G8YBXBX090SSZC3Y60H2` Fix Review question handles collide across runs started together: reviewQuestionHandle (internal/api/run_review.go ~927) builds review-<rid[:8]>-r<round>-<finding>, but a run ID is a ULID whose first 8 characters are only its start time (~35 s granularity), so plan runs the daemon starts in the same poll share the prefix. Evidence: plan runs of change sets 01M6VY7G6HPJ (dispatch) and 01M6VY7GRV7B (scheduling) both added question review-01m6vy7e-r3-f1 (nodes 01M6VZP9EAZX and 01M6VZK56QGX); rhiza health warns 'in-flight: model: duplicate handle', and when both change sets land rhiza check will fail on main. Use the run ID's random part (or the full ID, or the change set) in the handle, and when a handle is taken pick the next free one; add a test with two runs whose IDs share their first 8 characters. Existing duplicates: retitle one of the two when it lands (not unstuck yet).
