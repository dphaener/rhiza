# landed-green-retry-plan

source · `01M6XXH70ABSQH0W636M5GWRKD`

## quotes

```json
[]
```

## raw_sha256

d7adf395dbb40f05d660aaaf535031c82edd5fd974d6170cb9f6f7636bd1d0eb

## session

47f5cdc8-5adc-491d-a6cf-01a1bbd2cfd6

## summary

Autonomous plan run for the bug where ship-change's landed-green failed on pending CI (nebari PR #83). The planner found the cause: rhiza run resume --retry restarts a failed run at the failed command step (landed-green), skipping the wait step (wait-landed) that polls until checks settle. The change set was placed under m27-nebari-on-rhiza, like the plan-flake fix. One question (how a retry should wait for CI: a retry_at field, an engine rule inferring the wait, or green steps that poll) went to the arbiter, who decided retry-reopens-at-the-wait: a generic engine rule that reopens a retry at the wait step whose on.pass is the failed step. No owner words: the run was autonomous.

## title

Landed-green retry plan, autonomous

## Links

### In

- derived-from [retry-reopens-at-the-wait](../decision/retry-reopens-at-the-wait.md)
- derived-from [retry-reopens-at-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/retry-reopens-at-wait.md)

## History

- `01M6XX4R8Q2PCHDTQS80XFSG9J` Fix Ship-change's landed-green step fails instead of waiting when the PR's CI has not finished: nebari run 01M6WFRNHPRB (change set 01M6W90N0DB7, PR #83) was retried after a fix pushed head 5c49fea, and landed-green failed at once with 'PR #83 at 5c49fea: not settled: 0 of 8 checks finished', ending the run. A pending check is not a failure: landed-green (like wait-pr and rebased-green) should poll until the checks settle or a timeout passes, and fail only on a real red. Unstuck by waiting for CI by hand and resuming with --retry.
