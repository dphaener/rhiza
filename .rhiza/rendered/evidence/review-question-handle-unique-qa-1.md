# review-question-handle-unique-qa-1

evidence · `01M6WHA6QFSNWS04AJ0MDXJXCZ` · actor [claude-code-cb92742b](../actor/claude-code-cb92742b.md)

## result

PASS

## statement

QA by hand with the built binary (go build -trimpath ./cmd/rhiza) in a scratch Rhiza root with real wall-clock run IDs, no test clock. Two hold runs started in parallel got IDs 01M6WH7K22GY0EEKR47M709Y4N and 01M6WH7K5Y34JW19BVKB5XTRST, sharing the 8-character prefix the old handle used. With review-01m6wh7k22gy0eekr47m709y4n-r3-f1 pre-taken by an unrelated question, three findings rounds plus a confirmed f1 verdict on each run: review-check on run 1 added review-01m6wh7k22gy0eekr47m709y4n-r3-f1-2 and journalled that handle in its asked entry; on run 2 it added review-01m6wh7k5y34jw19bvkb5xtrst-r3-f1. A repeated review-check on each reported added false and left one asked entry per run; rhiza health showed no duplicate handle. An autonomous run with its full-ID handle pre-taken: rhiza run gate failed before review-check (no question left for f1) and passed after it, finding the -2 question by the journal. Scoped tests (go test ./internal/api -run Review; testscripts run-review-check, -handles, -landed, run-gate, run-arbitrate-asked) pass.

## Links

### Out

- verifies [review-question-handle-unique](../capability/rhiza/workflows/built-in-workflows/plan-workflow/review-question-handle-unique.md)

## History

- `01M6W7G8YBXBX090SSZC3Y60H2` Fix Review question handles collide across runs started together: reviewQuestionHandle (internal/api/run_review.go ~927) builds review-<rid[:8]>-r<round>-<finding>, but a run ID is a ULID whose first 8 characters are only its start time (~35 s granularity), so plan runs the daemon starts in the same poll share the prefix. Evidence: plan runs of change sets 01M6VY7G6HPJ (dispatch) and 01M6VY7GRV7B (scheduling) both added question review-01m6vy7e-r3-f1 (nodes 01M6VZP9EAZX and 01M6VZK56QGX); rhiza health warns 'in-flight: model: duplicate handle', and when both change sets land rhiza check will fail on main. Use the run ID's random part (or the full ID, or the change set) in the handle, and when a handle is taken pick the next free one; add a test with two runs whose IDs share their first 8 characters. Existing duplicates: retitle one of the two when it lands (not unstuck yet).
