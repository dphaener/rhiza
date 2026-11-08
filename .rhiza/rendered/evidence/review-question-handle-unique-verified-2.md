# review-question-handle-unique-verified-2

evidence · `01M6WJW3XM78MYZ7FWTD8WFNQ7` · actor [claude-code-a75fd010](../actor/claude-code-a75fd010.md)

## result

PASS

## statement

Verified at 7de305eb (attempt 2, after compliance finding f1). go test -trimpath -short ./internal/api -run 'Review|Gate|RunReview' -v: all 17 tests PASS, including TestReviewQuestionHandleFullID (two IDs sharing their first 8 chars get different handles carrying the full lower-cased ID) and TestReviewAskedFor, which now calls gateReview with a round-3 run whose journalled asked entry is review-<id>-r3-f1-2. The gate takes it as the question left, and it rejects another run's handle. (The package exit was FAIL only from the root-state guard: the root daemon changed ~/code/.rhiza-root/usage.live.json during the run. The full package hit a 10m timeout only because TestInboxGateQuestionsBudget blocked on the shared budgetLock under load 17.) Mutation check in a scratch copy: reverting gateReview to the old asked[reviewQuestionHandle(...)] lookup makes TestReviewAskedFor fail with 'confirmed finding f1 with no question left for it', so f1 is fixed. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-review-check.*|run-gate.*)$': ok. That covers run-review-check-handles: two runs started in the same second share an 8-char ID prefix; run 1's full-ID handle is taken by an unrelated question; review-check gives run 1 review-<id>-r3-f1-2, journals it in the asked entry, and gives run 2 its own base handle; a second review-check returns added=false with one asked entry. Diff read: reviewQuestionHandle uses the full ID; review-check finds journalled entries via reviewAskedFor, otherwise uniqueHandle; gateReview uses reviewAskedFor. Remaining nit f2 (a stale comment naming reviewQuestionHandleFor, and gateReview's asked map is now unused) does not affect acceptance.

## Links

### Out

- verifies [review-question-handle-unique](../capability/rhiza/workflows/built-in-workflows/plan-workflow/review-question-handle-unique.md)

## History

- `01M6W7G8YBXBX090SSZC3Y60H2` Fix Review question handles collide across runs started together: reviewQuestionHandle (internal/api/run_review.go ~927) builds review-<rid[:8]>-r<round>-<finding>, but a run ID is a ULID whose first 8 characters are only its start time (~35 s granularity), so plan runs the daemon starts in the same poll share the prefix. Evidence: plan runs of change sets 01M6VY7G6HPJ (dispatch) and 01M6VY7GRV7B (scheduling) both added question review-01m6vy7e-r3-f1 (nodes 01M6VZP9EAZX and 01M6VZK56QGX); rhiza health warns 'in-flight: model: duplicate handle', and when both change sets land rhiza check will fail on main. Use the run ID's random part (or the full ID, or the change set) in the handle, and when a handle is taken pick the next free one; add a test with two runs whose IDs share their first 8 characters. Existing duplicates: retitle one of the two when it lands (not unstuck yet).
