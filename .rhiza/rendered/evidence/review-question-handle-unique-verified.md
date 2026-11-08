# review-question-handle-unique-verified

evidence · `01M6WGZZN4Z8QH5JXJGDFPR5QQ` · actor [claude-code-a33034aa](../actor/claude-code-a33034aa.md)

## result

PASS

## statement

Verified at e58a2be7 by claude-code-a33034aa. Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-review-check-handles|run-review-check|run-review-check-landed|run-gate|run-arbitrate-asked)$': ok. In run-review-check-handles' -testwork dir the two runs were 01M5PWGKBGYFZMTH8Y8AF1G8GN and 01M5PWGKBGMVZ1H5AXG5VQ0BY9 (same first 10 chars, so the old 8-char handle would collide; the script asserts the shared 8-char prefix), both reached round 3 with f1 confirmed, and their journals hold asked entries with handles review-01m5pwgkbgyfzmth8y8af1g8gn-r3-f1-2 (run 1, whose full-ID handle an unrelated question held) and review-01m5pwgkbgmvz1h5axg5vq0by9-r3-f1 (run 2); a second review-check added nothing and kept one asked entry. go test -trimpath ./internal/api -run 'Review|Gate|Asked' (from a scratch worktree outside ~/code): ok, including TestReviewQuestionHandleFullID and TestReviewAskedFor (gate finds the -2 question via the journal). Read the diff: the gate and repeated checks use reviewAskedFor on journalled entries, not the rebuilt handle. Note: the full internal/api package timed out at 600s under VM load ~8.5, and a run under ~/code tripped the enclosing-root state guard; neither touches this change. Nit: a comment cites a nonexistent reviewQuestionHandleFor.

## Links

### Out

- verifies [review-question-handle-unique](../capability/rhiza/workflows/built-in-workflows/plan-workflow/review-question-handle-unique.md)

## History

- `01M6W7G8YBXBX090SSZC3Y60H2` Fix Review question handles collide across runs started together: reviewQuestionHandle (internal/api/run_review.go ~927) builds review-<rid[:8]>-r<round>-<finding>, but a run ID is a ULID whose first 8 characters are only its start time (~35 s granularity), so plan runs the daemon starts in the same poll share the prefix. Evidence: plan runs of change sets 01M6VY7G6HPJ (dispatch) and 01M6VY7GRV7B (scheduling) both added question review-01m6vy7e-r3-f1 (nodes 01M6VZP9EAZX and 01M6VZK56QGX); rhiza health warns 'in-flight: model: duplicate handle', and when both change sets land rhiza check will fail on main. Use the run ID's random part (or the full ID, or the change set) in the handle, and when a handle is taken pick the next free one; add a test with two runs whose IDs share their first 8 characters. Existing duplicates: retitle one of the two when it lands (not unstuck yet).
