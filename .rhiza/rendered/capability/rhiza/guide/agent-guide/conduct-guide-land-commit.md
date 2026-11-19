# conduct-guide-land-commit

capability · `01M7RJTNZVJ7N3FJ1M9N4G3PN3` · parent [agent-guide](../agent-guide.md)

## acceptance

go test -trimpath ./internal/guide -run 'TestConduct|TestConductHandback|TestBuildReviewCommandBudget' passes: TestConduct asserts conduct.md's Land phase names rhiza land --commit and the Change: and Actor: trailer rule; TestConductHandback asserts conduct-handback.md's closing land names rhiza land --commit; (failure case) the old 'rhiza land, then rhiza render, and commit what they change' wording is asserted absent from both topics

## files

internal/guide/text/conduct.md, internal/guide/text/conduct-handback.md, internal/guide/conduct_test.go, internal/guide/conduct_handback_test.go

## must_not

Change other roles' guide topics (builder, rhiza, change-sets, planner); add a phase to the conduct guide

## statement

The conduct guide's Land phase and conduct-handback's closing land tell the conductor to land, render and commit with rhiza land --commit, and say every other commit the conductor makes carries Change: and Actor: trailers, so none is journalled as a hand commit

## Links

### Out

- depends-on [land-commit](../../change-sets/change-lifecycle/land-commit.md)
- decided-by [conductor-land-commit-trailers](../../../../decision/conductor-land-commit-trailers.md)

### In

- verifies [conduct-guide-land-commit-verified](../../../../evidence/conduct-guide-land-commit-verified.md)
- verifies [conduct-guide-land-commit-qa-1](../../../../evidence/conduct-guide-land-commit-qa-1.md)

## History

- `01M7RJM81GCCYXK77G7HMKA0FR` Fix False hand-commit intervention: run 01M7RHR8NEJTECVZ5FYFDJ0H18 (change set 01M7RHC6VMMBXVWS3BKKFF6F8N) journals 'Interventions: 1 (hand-commit 1)' for commit 351cf62f6 'Land and render the run-status flake fix', which its own conductor (session 9239a99d, actor claude-code-9239a99d) made after rhiza land and rhiza render. The commit has no Change:/Actor:/Leaf:/Review: trailer, so handCommits (internal/api/run_handcommit.go) counts it as a hand commit and the daemon's clean streak breaks with no human involved. Earlier conductors happened to add trailers; nothing makes them. Make the conductor's own land/render commit always recognised (e.g. rhiza land/render commit with the run's trailers, or the conduct guide and a check requiring Change: and Actor: trailers on every conductor commit), and add a test where a conductor's land commit is not journalled as an intervention. This blocks milestone 40's measure (10 consecutive clean runs).
