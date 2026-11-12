# ship-check-skipped-jobs-green-verified

evidence · `01M76Q3ZX83T1J0SCX9WVXNHJS` · actor [claude-code-4ea94eea](../actor/claude-code-4ea94eea.md)

## result

PASS

## statement

Verifier at f87db76e: go test -trimpath ./internal/api -run 'ChecksVerdict|ShipRerunFailed|WatchdogCIRerun' -v passes (TestChecksVerdict/ModelOnlyPR: SKIPPED lint, scripts, api, rest, budgets plus SUCCESS changes, model, check read settled and green with no failed names; TestShipRerunFailedSkippedJobs: ship.check on that head is settled, green and passed with no Failed or Reran entries and no gh run rerun call; the new TestWatchdogCIRerun case feeds a head with the Go jobs SKIPPED and changes, model, check SUCCESS and asserts no rerun and that rerunSpent is false). Mutation check in a scratch copy: dropping SKIPPED from passingConclusions fails all three (ModelOnlyPR verdict red, ship check calls gh run rerun 42 --failed, watchdog 'skipped jobs reran [7]'). Note: the rerunSpent assertion uses Attempts 1 with CIReruns 1, so on its own it would pass even if SKIPPED counted as failed; the poll assertion in the same case does catch that. The diff only adds tests and the builder's actor record.

## Links

### Out

- verifies [ship-check-skipped-jobs-green](../capability/rhiza/orchestration/shipping/ship-check-skipped-jobs-green.md)

## History

- `01M7684F6CEBCS33NW5DDK1BVP` CI on a PR that changes only .rhiza/ (plan PRs, records, accepts, renders) runs just the model checks, not the Go suite, so a plan PR's CI takes a minute or two instead of 11-18. Owner, 2026-11-11, on the hour a milestone plan run spends after its interview: "an hour is also absurd for such a simple task", then "yes file it" for this request. Today .github/workflows/ci.yml runs lint, scripts (sharded), api (sharded), rest and budgets jobs on every PR, whatever it touches; plan PRs #435, #436 and #437 each took 15-18 minutes from interview to merge, mostly CI. Wanted: a first job decides from the PR's changed paths whether any file outside .rhiza/ changed (a paths filter, e.g. dorny/paths-filter or git diff against the merge base); when only .rhiza/ changed, CI builds rhiza and runs make self-check (rhiza check with this commit's build, including stale .rhiza/rendered) and skips the Go test jobs; any other change runs everything as today. The required status checks must still report on a model-only PR (skipped jobs count as passing, or one aggregate required job), so branch protection and rhiza ship check --settled treat it as green; check ship check and the watchdog's CI rerun still see the right checks. main's own CI (push to main) keeps running everything. Keep make check the single local entry point. Out of scope: the general CI speed-up (request 01M6Q58QG7T8V9F602YX6J2XPZ).
