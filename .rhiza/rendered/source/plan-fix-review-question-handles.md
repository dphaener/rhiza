# plan-fix-review-question-handles

source · `01M6W7Q7V0T5NYNNDYC1DHZ13Q`

## quotes

```json
[]
```

## raw_sha256

9011dbad93b2e4478a2a431df03852e4a6f6b34a6f0f2964d5748d20b9467cc3

## session

c97ce25f-459e-4c9d-8e24-01a1b8787167

## summary

Autonomous plan run for change set 01M6W7G8 (Fix Review question handles). No owner took part. The planner read internal/api/run_review.go: reviewQuestionHandle builds review-<rid[:8]>-r<round>-<finding>. The first 8 characters of a ULID are only its start time, so plan runs the daemon starts in the same poll share them; runs of 01M6VY7G6HPJ and 01M6VY7GRV7B both added review-01m6vy7e-r3-f1. The handle is also the identity key: reviewQuestionNode reuses a question by handle across change sets, and the asked journal and the autonomous accept gate rebuild it. The planner raised one question: what the handle carries and what happens when it is taken. The arbiter answered with decision review-question-handle-full-run-id (assumed, would have escalated): the full lower-cased run ID; when that handle is taken by another node, take the first free one with uniqueHandle and journal the handle actually used; later checks match by the journal entry. Planned: one leaf, review-question-handle-unique, under plan-workflow, serving milestone-29, whose plan runs hit the collision. Not planned here: the two existing review-01m6vy7e-r3-f1 nodes belong to other unlanded change sets; the intent says to retitle one of them when it lands.

## title

Plan of Fix Review question handles (autonomous)

## Links

### In

- derived-from [review-question-handle-full-run-id](../decision/review-question-handle-full-run-id.md)

## History

- `01M6W7G8YBXBX090SSZC3Y60H2` Fix Review question handles collide across runs started together: reviewQuestionHandle (internal/api/run_review.go ~927) builds review-<rid[:8]>-r<round>-<finding>, but a run ID is a ULID whose first 8 characters are only its start time (~35 s granularity), so plan runs the daemon starts in the same poll share the prefix. Evidence: plan runs of change sets 01M6VY7G6HPJ (dispatch) and 01M6VY7GRV7B (scheduling) both added question review-01m6vy7e-r3-f1 (nodes 01M6VZP9EAZX and 01M6VZK56QGX); rhiza health warns 'in-flight: model: duplicate handle', and when both change sets land rhiza check will fail on main. Use the run ID's random part (or the full ID, or the change set) in the handle, and when a handle is taken pick the next free one; add a test with two runs whose IDs share their first 8 characters. Existing duplicates: retitle one of the two when it lands (not unstuck yet).
