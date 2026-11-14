# review-01m7ah6rayrzax6z6r63q0m2qk-r6-f2

question · `01M7BZ7EB4840AK4XWZXBGCP7Y`

## statement

Plan review round 6 left finding f2 (contradiction on review-check-declared-checks): Already on main: conductorReviewCheck in internal/api/run_review.go (01M79KGZ, leaf critic-on-major) runs when a step declares checks. It reads each declared check's latest round and fails on standing blocker or major findings. Main also fails when a declared check recorded no round, but this leaf says review-check passes in every other case. The plan re-specifies landed behaviour and contradicts the code on missing rounds. The critic confirmed it: Verified: origin/main internal/api/run_review.go:1235 conductorReviewCheck runs for a step that declares checks, and at :1267 it fails when a declared check recorded no round. The leaf says review-check passes in every case other than a standing blocker or major, which contradicts main on the missing-round case and re-specifies landed work. Major stands. How should it be settled?

## Links

### Out

- blocks [review-check-declared-checks](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-check-declared-checks.md)

### In

- answers [review-check-declared-checks-follows-main](../decision/review-check-declared-checks-follows-main.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
