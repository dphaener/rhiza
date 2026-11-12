# ci-model-only-pr

capability · `01M768E0NH089HZ4WRR0SBBDAB` · parent [self-check](../self-check.md)

## acceptance

Go tests in internal/ciworkflow (go test -trimpath ./internal/ciworkflow/...) parse .github/workflows/ci.yml and assert: a `changes` job with an output code, computing it from git diff --name-only against the merge base on pull_request and setting it true on workflow_dispatch; lint, scripts, api, rest and budgets each needing changes with if on code == 'true'; a `model` job needing changes with if on code == 'false' that runs make self-check and `go test -trimpath ./internal/op -run TestEveryFieldInRepoLogIsKnown` (so a model-only PR whose log carries a field op.FieldKnown rejects still fails CI, as it does today through the api/rest jobs); check needing changes, model and the five Go jobs, if: always(), failing unless changes succeeded and every other result is success or skipped, and failing when none of model and the five Go jobs succeeded (an all-skipped run fails check); the existing triggers, concurrency, matrices and cache checks still pass. make check is unchanged. The change set's own PR touches Go code so runs everything; the next model-only PR's CI completes in a few minutes with check green, read from gh run view.

## statement

On a pull request whose changed paths all lie under .rhiza/, CI runs only a `model` job (setup-go, make self-check, then `go test -trimpath ./internal/op -run TestEveryFieldInRepoLogIsKnown`, the one Go test that reads the repo's own .rhiza/log and fails on a node field this commit's code does not know, which neither fold nor rhiza check validates) and skips lint, scripts, api, rest and budgets, decided by a first `changes` job using plain git diff against the merge base (code=true when any path is outside .rhiza/, always true on workflow_dispatch); any other PR and every dispatched run on main runs everything as today, and the aggregate `check` job, needing changes, model and the five Go jobs with if: always(), fails unless changes succeeded and passes only when every other job is success or skipped with at least one of them succeeded.

## Links

### Out

- decided-by [decision-model-only-ci](../../../../decision/decision-model-only-ci.md)

### In

- verifies [ci-model-only-pr-verified](../../../../evidence/ci-model-only-pr-verified.md)
- verifies [ci-model-only-pr-qa-1](../../../../evidence/ci-model-only-pr-qa-1.md)
- verifies [ci-model-only-pr-verified-3](../../../../evidence/ci-model-only-pr-verified-3.md)
- verifies [pr-458-ci-on-a-pr-that-changes-only](../../../../evidence/pr-458-ci-on-a-pr-that-changes-only.md)

## History

- `01M7684F6CEBCS33NW5DDK1BVP` CI on a PR that changes only .rhiza/ (plan PRs, records, accepts, renders) runs just the model checks, not the Go suite, so a plan PR's CI takes a minute or two instead of 11-18. Owner, 2026-11-11, on the hour a milestone plan run spends after its interview: "an hour is also absurd for such a simple task", then "yes file it" for this request. Today .github/workflows/ci.yml runs lint, scripts (sharded), api (sharded), rest and budgets jobs on every PR, whatever it touches; plan PRs #435, #436 and #437 each took 15-18 minutes from interview to merge, mostly CI. Wanted: a first job decides from the PR's changed paths whether any file outside .rhiza/ changed (a paths filter, e.g. dorny/paths-filter or git diff against the merge base); when only .rhiza/ changed, CI builds rhiza and runs make self-check (rhiza check with this commit's build, including stale .rhiza/rendered) and skips the Go test jobs; any other change runs everything as today. The required status checks must still report on a model-only PR (skipped jobs count as passing, or one aggregate required job), so branch protection and rhiza ship check --settled treat it as green; check ship check and the watchdog's CI rerun still see the right checks. main's own CI (push to main) keeps running everything. Keep make check the single local entry point. Out of scope: the general CI speed-up (request 01M6Q58QG7T8V9F602YX6J2XPZ).
