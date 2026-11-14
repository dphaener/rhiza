# critic-on-major

capability · `01M79KPDBW0RMPYBW73RTEBNX3` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

go test -trimpath ./internal/api -run TestReviewCheckConductor shows, in a run whose conduct step declares plan, code and qa and whose checks rounds are journaled under conduct, with review-check run while the run's current step is the later gate command step that declares no checks: a round of minor findings passing with no critic verdict, a blocker with no verdict and a confirmed major failing with exit non-zero and naming them, a blocker downgraded to minor and a dropped major passing, only each check's latest round counted, and no 'has no review round yet' error from the plan-review path. Touches internal/api/run_review.go, internal/api/run_review_test.go.

## statement

rhiza run review-check works in a run whose workflow has a step that declares checks, from that step or any later step (the conductor default's gate, a command step after conduct that declares none): it takes the declared checks from the run workflow's step that declares them, not from the current step, reads the run's journal for the rounds recorded under that step's name, and, for each declared check, counts its latest round's blocker and major findings that no critic verdict downgraded or dropped as open, failing (exit non-zero, listing them) while any is open; a round with only minor or nit findings passes with no critic. It never falls back to the run-level plan-review rounds in such a run. The count comes from the journal, never from the conductor's word.

## Links

### Out

- depends-on [review-three-checks](review-three-checks.md)
- decided-by [conductor-three-checks](../../../../decision/conductor-three-checks.md)
- decided-by [conductor-checks-gate-pick](../../../../decision/conductor-checks-gate-pick.md)

### In

- verifies [critic-on-major-verified](../../../../evidence/critic-on-major-verified.md)
- verifies [critic-on-major-qa-1](../../../../evidence/critic-on-major-qa-1.md)
- verifies [critic-on-major-verified-2](../../../../evidence/critic-on-major-verified-2.md)
- verifies [critic-on-major-qa-2](../../../../evidence/critic-on-major-qa-2.md)
- verifies [pr-575-a-rhiza-conduct-skill-conducts](../../../../evidence/pr-575-a-rhiza-conduct-skill-conducts.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
