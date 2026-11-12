# ship-check-skipped-jobs-green

capability · `01M768E0S5MMHEP3CNXNJYBYPA` · parent [shipping](../shipping.md)

## acceptance

Go tests in internal/api (go test -trimpath ./internal/api -run 'ChecksVerdict|ShipRerunFailed|WatchdogCIRerun') feed checksVerdict a rollup of SKIPPED lint, scripts, api, rest and budgets plus SUCCESS changes, model and check and get settled and green with no failed names; run shipRerunFailed on such a head and assert it issues no gh run rerun; and, in a TestWatchdogCIRerun case in internal/api/daemon_watchdog_test.go, feed watchdogCIRerun a PR whose head has the Go jobs SKIPPED and changes, model and check SUCCESS and assert it issues no gh run rerun and that rerunSpent does not count the rerun spent.

## statement

ship check reads a model-only PR's CI as settled and green: a check rollup whose skipped jobs conclude SKIPPED and whose `check` and `model` jobs conclude SUCCESS is settled and green, and the CI rerun reruns nothing for SKIPPED jobs.

## Links

### Out

- decided-by [decision-model-only-ci](../../../../decision/decision-model-only-ci.md)

### In

- verifies [ship-check-skipped-jobs-green-verified](../../../../evidence/ship-check-skipped-jobs-green-verified.md)
- verifies [ship-check-skipped-jobs-green-qa-1](../../../../evidence/ship-check-skipped-jobs-green-qa-1.md)
- verifies [pr-458-ci-on-a-pr-that-changes-only](../../../../evidence/pr-458-ci-on-a-pr-that-changes-only.md)

## History

- `01M7684F6CEBCS33NW5DDK1BVP` CI on a PR that changes only .rhiza/ (plan PRs, records, accepts, renders) runs just the model checks, not the Go suite, so a plan PR's CI takes a minute or two instead of 11-18. Owner, 2026-11-11, on the hour a milestone plan run spends after its interview: "an hour is also absurd for such a simple task", then "yes file it" for this request. Today .github/workflows/ci.yml runs lint, scripts (sharded), api (sharded), rest and budgets jobs on every PR, whatever it touches; plan PRs #435, #436 and #437 each took 15-18 minutes from interview to merge, mostly CI. Wanted: a first job decides from the PR's changed paths whether any file outside .rhiza/ changed (a paths filter, e.g. dorny/paths-filter or git diff against the merge base); when only .rhiza/ changed, CI builds rhiza and runs make self-check (rhiza check with this commit's build, including stale .rhiza/rendered) and skips the Go test jobs; any other change runs everything as today. The required status checks must still report on a model-only PR (skipped jobs count as passing, or one aggregate required job), so branch protection and rhiza ship check --settled treat it as green; check ship check and the watchdog's CI rerun still see the right checks. main's own CI (push to main) keeps running everything. Keep make check the single local entry point. Out of scope: the general CI speed-up (request 01M6Q58QG7T8V9F602YX6J2XPZ).
