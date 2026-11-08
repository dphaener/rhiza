# review-question-handle-full-run-id

decision · `01M6W7P61Z5C2BECZF7K9RV4KF`

## reasoning

The change set's intent offers three carriers (the ULID's random part, the full ID, or the change set) and asks for a next-free fallback plus a test with two runs sharing their first 8 characters. Reading internal/api/run_review.go: the handle is not just a label but the identity key: reviewQuestionNode looks a question up by handle across change sets so a repeated check does not add it twice, and both the asked journal and the autonomous accept gate (~line 1494) recompute reviewQuestionHandle(rid, round, finding). So the handle must stay a pure function of (run, round, finding) and be unique per run. The full run ID is unique by construction, needs no reasoning about which ULID characters are random, and reads back directly as the ID 'rhiza run status' takes; the change set alone would not separate two runs of the same change set. The repo already has uniqueHandle (workspace.go) for taken handles, used by sessions, questions and ship close; reusing it covers handles taken by unrelated nodes or by the existing review-01m6vy7e-* duplicates, and journalling the used handle keeps idempotence and the gate correct when the fallback fires.

## statement

A review question's handle carries the run's full ID, lower-cased: review-<run ID>-r<round>-<finding> (about 40 characters, under the 64 limit). It stays deterministic, so a re-run review-check finds and reuses its own question rather than adding another. When that handle is already taken by a node that is not this run's question for this finding, review-check takes the first free one with uniqueHandle (base-2, base-3, ...) and records the handle it actually used in the asked journal entry; the run's later checks (the accept gate's 'question left for each confirmed finding') match its questions by the journalled entry, not by recomputing the handle.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-c97ce25f-1](../question/question-c97ce25f-1.md)
- derived-from [plan-fix-review-question-handles](../source/plan-fix-review-question-handles.md)

### In

- decided-by [review-question-handle-unique](../capability/rhiza/workflows/built-in-workflows/plan-workflow/review-question-handle-unique.md)

## History

- `01M6W7G8YBXBX090SSZC3Y60H2` Fix Review question handles collide across runs started together: reviewQuestionHandle (internal/api/run_review.go ~927) builds review-<rid[:8]>-r<round>-<finding>, but a run ID is a ULID whose first 8 characters are only its start time (~35 s granularity), so plan runs the daemon starts in the same poll share the prefix. Evidence: plan runs of change sets 01M6VY7G6HPJ (dispatch) and 01M6VY7GRV7B (scheduling) both added question review-01m6vy7e-r3-f1 (nodes 01M6VZP9EAZX and 01M6VZK56QGX); rhiza health warns 'in-flight: model: duplicate handle', and when both change sets land rhiza check will fail on main. Use the run ID's random part (or the full ID, or the change set) in the handle, and when a handle is taken pick the next free one; add a test with two runs whose IDs share their first 8 characters. Existing duplicates: retitle one of the two when it lands (not unstuck yet).
