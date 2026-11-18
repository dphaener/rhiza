# open-pr-review-when-declared-verified

evidence · `01M7RC9X0VPY8EVKGZWHQPGT9Y` · actor [claude-code-f827b748](../actor/claude-code-f827b748.md)

## result

PASS

## statement

At 90af1739a, under the capped scope: go test -trimpath -p 2 ./internal/api -run 'TestShipOpenPR|TestReviewJournaled|TestReviewDeclared' exit 0 (TestShipOpenPRReviewless: myteam-flow from docs/examples/byo/workflow.toml and a review-less workflow open the PR with no review; TestShipOpenPRRefusals: review run with none journaled refused naming why, call outside a run refused; TestReviewDeclared true for ship-change, ship-change-classic, review, bugfix, plan runs, missing/unparsable definition, false for myteam-flow and a review-less workflow named review); go test -trimpath -p 2 ./internal/workflow exit 0 (TestDeclaresReview by kinds/[[check]] tables/included workflows, never names; a skipped included checks step declares none); RHIZA_*/SHIP_*/REX_* unset go test ./cmd/rhiza -run 'TestScripts/(ship-open-pr.*|run-ship-change)$' exit 0; go vet exit 0; rhiza ship generate and rhiza render exit 0 with no changes. Diff read: DeclaresReview checks Kind==checks or Conducts(); guard kept for plan runs and every built-in/repo workflow with a review (docs, ship-rhiza included); no workflow file key added, docs/examples/byo untouched, no TODOs.

## Links

### Out

- verifies [open-pr-review-when-declared](../capability/ship-through-rhiza/open-pr-review-when-declared.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
