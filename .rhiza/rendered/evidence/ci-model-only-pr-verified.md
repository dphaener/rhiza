# ci-model-only-pr-verified

evidence · `01M76Q13QPKGZE13M6PMJ69CHG` · actor [claude-code-255a77a1](../actor/claude-code-255a77a1.md)

## result

PASS

## statement

Verifier at 99bb0522f: go test -trimpath -count=1 ./internal/ciworkflow/... passes all 8 tests (incl. new TestCIChangesDecidesFromGitDiff, TestCIGoJobsRunOnlyOnCodeChanges, TestCIModelJobRunsOnlyOnModelChanges, updated TestCICheckAggregatesEveryJob and the existing trigger/concurrency/matrix/checkout/cache tests). Diff read: ci.yml adds a changes job (output code from steps.diff; code=true unless event is pull_request, else git merge-base + git diff --name-only filtering out .rhiza/), a model job (needs changes, if code == 'false', setup-go cache: true, make self-check, go test -trimpath ./internal/op -run TestEveryFieldInRepoLogIsKnown), needs/if code == 'true' on lint, scripts, api, rest, budgets, and check needing all seven with if: always(). Makefile untouched. Extracted and ran the check script: model-only results and code results exit 0; all-skipped, changes failure, any failure or cancelled exit 1. Extracted and ran the changes script: on HEAD vs parent code=true, on model-only commit f7a79b045 vs parent code=false, on workflow_dispatch code=true. ci.yml parses as YAML; TestEveryFieldInRepoLogIsKnown passes. Not verifiable pre-merge: the next model-only PR's live CI timing via gh run view.

## Links

### Out

- verifies [ci-model-only-pr](../capability/rhiza/integrity/self-check/ci-model-only-pr.md)

## History

- `01M7684F6CEBCS33NW5DDK1BVP` CI on a PR that changes only .rhiza/ (plan PRs, records, accepts, renders) runs just the model checks, not the Go suite, so a plan PR's CI takes a minute or two instead of 11-18. Owner, 2026-11-11, on the hour a milestone plan run spends after its interview: "an hour is also absurd for such a simple task", then "yes file it" for this request. Today .github/workflows/ci.yml runs lint, scripts (sharded), api (sharded), rest and budgets jobs on every PR, whatever it touches; plan PRs #435, #436 and #437 each took 15-18 minutes from interview to merge, mostly CI. Wanted: a first job decides from the PR's changed paths whether any file outside .rhiza/ changed (a paths filter, e.g. dorny/paths-filter or git diff against the merge base); when only .rhiza/ changed, CI builds rhiza and runs make self-check (rhiza check with this commit's build, including stale .rhiza/rendered) and skips the Go test jobs; any other change runs everything as today. The required status checks must still report on a model-only PR (skipped jobs count as passing, or one aggregate required job), so branch protection and rhiza ship check --settled treat it as green; check ship check and the watchdog's CI rerun still see the right checks. main's own CI (push to main) keeps running everything. Keep make check the single local entry point. Out of scope: the general CI speed-up (request 01M6Q58QG7T8V9F602YX6J2XPZ).
