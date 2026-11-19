# conduct-guide-land-commit-verified

evidence · `01M7RKG68GYVJJXHEX2STJKX6J` · actor [claude-code-2dbd263e](../actor/claude-code-2dbd263e.md)

## result

PASS

## statement

Detached checkout at 1152041d3. Re-ran Checks: go test -trimpath ./internal/guide -run 'TestConduct|TestConductHandback|TestBuildReviewCommandBudget' exit 0 (TestConduct, TestConductHandback, TestBuildReviewCommandBudget PASS); go test ./internal/guide/... ok; rhiza ship generate exit 0, 0 changed; rhiza render exit 0. Diff: conduct.md Land phase now names rhiza land --commit (lands, renders, commits .rhiza with Change:/Actor: trailers) and says every other conductor commit carries Change: and Actor: trailers so none is journalled as a hand commit; conduct-handback.md closing land uses rhiza land --commit. TestConduct asserts both; oldLandWording asserted absent from both topics. Failure case checked by temporarily restoring the old 'commit what they change' wording in conduct-handback.md: TestConductHandback FAILED with 'still says', then reverted. Only the leaf's four files plus the builder's own actor record changed; no other guide topics touched, no phase added, no TODOs. Minor: rhiza render at the commit adds the actor's line to .rhiza/rendered/index.md, which the commit left out; the conductor's land --commit re-renders it.

## Links

### Out

- verifies [conduct-guide-land-commit](../capability/rhiza/guide/agent-guide/conduct-guide-land-commit.md)

## History

- `01M7RJM81GCCYXK77G7HMKA0FR` Fix False hand-commit intervention: run 01M7RHR8NEJTECVZ5FYFDJ0H18 (change set 01M7RHC6VMMBXVWS3BKKFF6F8N) journals 'Interventions: 1 (hand-commit 1)' for commit 351cf62f6 'Land and render the run-status flake fix', which its own conductor (session 9239a99d, actor claude-code-9239a99d) made after rhiza land and rhiza render. The commit has no Change:/Actor:/Leaf:/Review: trailer, so handCommits (internal/api/run_handcommit.go) counts it as a hand commit and the daemon's clean streak breaks with no human involved. Earlier conductors happened to add trailers; nothing makes them. Make the conductor's own land/render commit always recognised (e.g. rhiza land/render commit with the run's trailers, or the conduct guide and a check requiring Change: and Actor: trailers on every conductor commit), and add a test where a conductor's land commit is not journalled as an intervention. This blocks milestone 40's measure (10 consecutive clean runs).
