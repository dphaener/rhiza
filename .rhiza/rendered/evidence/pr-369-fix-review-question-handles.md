# pr-369-fix-review-question-handles

evidence · `01M6WKT98RK1ZX27BFSWGYA3SK` · actor [dphaener](../actor/dphaener.md)

## head

bd5d9447c5097bbaf3310a0349aa05f8fe1d0938

## pr

https://github.com/dphaener/rhiza/pull/369

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #369 (https://github.com/dphaener/rhiza/pull/369) passed CI at bd5d944

## status

verified

## Links

### Out

- verifies [review-question-handle-unique](../capability/rhiza/workflows/built-in-workflows/plan-workflow/review-question-handle-unique.md)

## History

- `01M6W7G8YBXBX090SSZC3Y60H2` Fix Review question handles collide across runs started together: reviewQuestionHandle (internal/api/run_review.go ~927) builds review-<rid[:8]>-r<round>-<finding>, but a run ID is a ULID whose first 8 characters are only its start time (~35 s granularity), so plan runs the daemon starts in the same poll share the prefix. Evidence: plan runs of change sets 01M6VY7G6HPJ (dispatch) and 01M6VY7GRV7B (scheduling) both added question review-01m6vy7e-r3-f1 (nodes 01M6VZP9EAZX and 01M6VZK56QGX); rhiza health warns 'in-flight: model: duplicate handle', and when both change sets land rhiza check will fail on main. Use the run ID's random part (or the full ID, or the change set) in the handle, and when a handle is taken pick the next free one; add a test with two runs whose IDs share their first 8 characters. Existing duplicates: retitle one of the two when it lands (not unstuck yet).
