# ship-head-build-step-qa-2

evidence · `01M77CPNGJGKFAB8V7S7B3488Y` · actor [claude-code-6f75672e](../actor/claude-code-6f75672e.md)

## result

FAIL

## statement

QA at 3bda3170 with a -trimpath build, in a scratch repo with a bare origin and a fake gh: on a pushed change set branch behind a main that had landed and rendered other change sets, ran rhiza ship rebase with SHIP_HEAD_RHIZA and SHIP_HEAD_BUILD naming a logging build script. Clean merge: the build ran after the merge (its tree held main's new change set log), then render --json ran through SHIP_HEAD_RHIZA, the merged head was pushed and rhiza check printed ok. Rendered-only conflict (branch rewrote rendered/index.md): the build ran during the merge (MERGE_HEAD set, main's logs present) before the settling render, and again before the final render; Settled 1 rendered file, pushed, check ok. Failing SHIP_HEAD_BUILD: exit 1 with the build's stderr, the branch and origin's branch kept their head, no merge left in progress, clean tree. go test -trimpath ./internal/api -run 'TestShipRebase|TestSessionRun|TestShipLand' passes (only the known ~/code .rhiza-root usage guard trips). But TestScripts/ship-land-branch-build now FAILS: it still expects ship land to go build rhiza when go.mod names Rhiza's module, the isRhizaRepo path this change set removed.

## Links

### Out

- verifies [ship-head-build-step](../capability/rhiza/orchestration/shipping/ship-head-build-step.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
