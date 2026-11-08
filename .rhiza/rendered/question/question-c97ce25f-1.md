# question-c97ce25f-1

question · `01M6W7M77J37T0KR2T3CR2B39Q`

## session

c97ce25f-459e-4c9d-8e24-01a1b8787167

## statement

Review question handles (run.review-check) collide because they use only the run ID's first 8 characters, which are its start time. What should the handle carry, and how should a taken handle be handled?

## Links

### In

- answers [review-question-handle-full-run-id](../decision/review-question-handle-full-run-id.md)

## History

- `01M6W7G8YBXBX090SSZC3Y60H2` Fix Review question handles collide across runs started together: reviewQuestionHandle (internal/api/run_review.go ~927) builds review-<rid[:8]>-r<round>-<finding>, but a run ID is a ULID whose first 8 characters are only its start time (~35 s granularity), so plan runs the daemon starts in the same poll share the prefix. Evidence: plan runs of change sets 01M6VY7G6HPJ (dispatch) and 01M6VY7GRV7B (scheduling) both added question review-01m6vy7e-r3-f1 (nodes 01M6VZP9EAZX and 01M6VZK56QGX); rhiza health warns 'in-flight: model: duplicate handle', and when both change sets land rhiza check will fail on main. Use the run ID's random part (or the full ID, or the change set) in the handle, and when a handle is taken pick the next free one; add a test with two runs whose IDs share their first 8 characters. Existing duplicates: retitle one of the two when it lands (not unstuck yet).
