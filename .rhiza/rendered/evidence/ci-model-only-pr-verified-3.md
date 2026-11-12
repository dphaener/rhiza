# ci-model-only-pr-verified-3

evidence · `01M76QVXQXF8H6622PMAGYNTS4` · actor [claude-code-ded57ce2](../actor/claude-code-ded57ce2.md)

## result

PASS

## statement

Verifier, attempt 3, at ff79bd5e5: go test -trimpath -count=1 ./internal/ciworkflow/... passes (TestCIChangesDecidesFromGitDiff, TestCIGoJobsRunOnlyOnCodeChanges, TestCIModelJobRunsOnlyOnModelChanges, TestCICheckAggregatesEveryJob plus the existing trigger, checkout/run-name, concurrency and jobs/matrices tests). Read .github/workflows/ci.yml: changes outputs code from git diff --name-only against git merge-base of the PR base and HEAD, code=true on non-pull_request events; lint, scripts, api, rest, budgets need changes with if code == 'true'; model needs changes with if code == 'false' and runs make self-check then go test -trimpath ./internal/op -run TestEveryFieldInRepoLogIsKnown (passes locally); check needs changes, model and the five Go jobs, if: always(), fails unless changes succeeded, fails on any result other than success/skipped, and fails when only changes succeeded. Makefile has no diff from the merge-base. The retry's commit 5cdf2698e amends ci-yml-parallel-jobs' statement and acceptance to match the new check job, closing the model finding. make self-check reports only stale-render fails for this run's own untracked actor records, refreshed at close-out. The clause about the next model-only PR's gh run view can only be read after merge.

## Links

### Out

- verifies [ci-model-only-pr](../capability/rhiza/integrity/self-check/ci-model-only-pr.md)

## History

- `01M7684F6CEBCS33NW5DDK1BVP` CI on a PR that changes only .rhiza/ (plan PRs, records, accepts, renders) runs just the model checks, not the Go suite, so a plan PR's CI takes a minute or two instead of 11-18. Owner, 2026-11-11, on the hour a milestone plan run spends after its interview: "an hour is also absurd for such a simple task", then "yes file it" for this request. Today .github/workflows/ci.yml runs lint, scripts (sharded), api (sharded), rest and budgets jobs on every PR, whatever it touches; plan PRs #435, #436 and #437 each took 15-18 minutes from interview to merge, mostly CI. Wanted: a first job decides from the PR's changed paths whether any file outside .rhiza/ changed (a paths filter, e.g. dorny/paths-filter or git diff against the merge base); when only .rhiza/ changed, CI builds rhiza and runs make self-check (rhiza check with this commit's build, including stale .rhiza/rendered) and skips the Go test jobs; any other change runs everything as today. The required status checks must still report on a model-only PR (skipped jobs count as passing, or one aggregate required job), so branch protection and rhiza ship check --settled treat it as green; check ship check and the watchdog's CI rerun still see the right checks. main's own CI (push to main) keeps running everything. Keep make check the single local entry point. Out of scope: the general CI speed-up (request 01M6Q58QG7T8V9F602YX6J2XPZ).
