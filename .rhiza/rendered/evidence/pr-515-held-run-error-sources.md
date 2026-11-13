# pr-515-held-run-error-sources

evidence · `01M799Q25699SZX31EFA4KJY6S` · actor [dphaener](../actor/dphaener.md)

## head

e260c16d540c6ca7dd0af4c3acc85189051fd022

## pr

https://github.com/dphaener/rhiza/pull/515

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #515 (https://github.com/dphaener/rhiza/pull/515) passed CI at e260c16

## status

verified

## Links

### Out

- verifies [ship-gh-transient-retry](../capability/rhiza/orchestration/shipping/ship-gh-transient-retry.md)
- verifies [ship-push-merges-own-origin](../capability/rhiza/orchestration/shipping/ship-push-merges-own-origin.md)
- verifies [land-conflict-to-fixer](../capability/rhiza/orchestration/shipping/land-conflict-to-fixer.md)
- verifies [actor-record-no-duplicate](../capability/rhiza/lineage/actors/actor-record-no-duplicate.md)
- verifies [rex-tab-launch-retry](../capability/rhiza/runtime/rex-sessions/rex-tab-launch/rex-tab-launch-retry.md)
- verifies [ship-merge-refused-stops](../capability/rhiza/orchestration/shipping/ship-merge-refused-stops.md)

## History

- `01M76XSQKQJ8P78585NDYPATM7` Remove the error sources that held ship-change runs: retry GitHub 5xx and GraphQL errors with backoff, stop once with GitHub's reason instead of looping on a refused merge, reconcile an origin branch that holds the run's own earlier pushes, and turn land's PR-head conflicts, duplicate actor handles and Rex tab failures into a retry or one clear stop. Evidence: 59 holds, 97.6h ALL; open-pr 13, gate-rebase-merge max_visits 9, land 6, code-check 4, Rex tab 4 (perf report section d, rec 3).
