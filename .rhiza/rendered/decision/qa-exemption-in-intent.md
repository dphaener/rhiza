# qa-exemption-in-intent

decision · `01M650RAVJZ1JHTPMKY658WE7B`

## reasoning

The question asks whether the intent or the decision should change. The decision fakes-only-in-build-and-review and leaf review-guide-fakes-only already carve out the qa check, internal/guide/text/review.md says every check except qa tests runs with fakes, and internal/guide/review_test.go pins that. Only the intent states the broader no-exception rule. Aligning the intent to the decision and code is the smaller change and leaves the agreed behaviour unchanged. A change set's intent lives in its open log entry, which no rhiza command rewrites and the append-only log forbids editing by hand, so this decision, not a new intent text, carries the exemption: wherever the intent and this decision differ, this decision governs.

## statement

The change set intent, whose text is fixed once the change set opens and cannot be amended, is read with the same qa exemption as the decision fakes-only-in-build-and-review: the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh, except the review qa check, which runs the built binary hands-on, and fakes alone do not count for it.

## wording

(assumed, would have escalated)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
