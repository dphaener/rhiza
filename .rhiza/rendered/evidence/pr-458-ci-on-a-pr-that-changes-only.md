# pr-458-ci-on-a-pr-that-changes-only

evidence · `01M76RPE2FDC57KKVGYDBK51E3` · actor [dphaener](../actor/dphaener.md)

## head

27cfbcb41501257c7b7ae1baafe4a26db279dbcf

## pr

https://github.com/dphaener/rhiza/pull/458

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #458 (https://github.com/dphaener/rhiza/pull/458) passed CI at 27cfbcb

## status

verified

## Links

### Out

- verifies [ci-model-only-pr](../capability/rhiza/integrity/self-check/ci-model-only-pr.md)
- verifies [ship-check-skipped-jobs-green](../capability/rhiza/orchestration/shipping/ship-check-skipped-jobs-green.md)
- verifies [ci-yml-parallel-jobs](../capability/rhiza/integrity/self-check/ci-yml-parallel-jobs.md)

## History

- `01M7684F6CEBCS33NW5DDK1BVP` CI on a PR that changes only .rhiza/ (plan PRs, records, accepts, renders) runs just the model checks, not the Go suite, so a plan PR's CI takes a minute or two instead of 11-18. Owner, 2026-11-11, on the hour a milestone plan run spends after its interview: "an hour is also absurd for such a simple task", then "yes file it" for this request. Today .github/workflows/ci.yml runs lint, scripts (sharded), api (sharded), rest and budgets jobs on every PR, whatever it touches; plan PRs #435, #436 and #437 each took 15-18 minutes from interview to merge, mostly CI. Wanted: a first job decides from the PR's changed paths whether any file outside .rhiza/ changed (a paths filter, e.g. dorny/paths-filter or git diff against the merge base); when only .rhiza/ changed, CI builds rhiza and runs make self-check (rhiza check with this commit's build, including stale .rhiza/rendered) and skips the Go test jobs; any other change runs everything as today. The required status checks must still report on a model-only PR (skipped jobs count as passing, or one aggregate required job), so branch protection and rhiza ship check --settled treat it as green; check ship check and the watchdog's CI rerun still see the right checks. main's own CI (push to main) keeps running everything. Keep make check the single local entry point. Out of scope: the general CI speed-up (request 01M6Q58QG7T8V9F602YX6J2XPZ).
