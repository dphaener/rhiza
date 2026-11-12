# ci-model-only-pr-qa-1

evidence · `01M76QBB464W3WN3CBMNA1KY5F` · actor [claude-code-d5314015](../actor/claude-code-d5314015.md)

## result

PASS

## statement

QA at 7f99063: extracted the changes and check job scripts from .github/workflows/ci.yml and ran them in a scratch git repo. changes: a PR touching only .rhiza/ gives code=false (also when the base advanced on main after branching, via merge-base); a PR also touching main.go gives code=true; workflow_dispatch and push give code=true. check: [changes+model success, Go jobs skipped] exits 0; [model skipped, Go jobs success] exits 0; all skipped but changes exits 1; a failed Go job, a failed model job, a cancelled job, or a failed changes job each exit 1. Job graph parsed with PyYAML: model needs changes if code == 'false', lint/scripts/api/rest/budgets need changes if code == 'true', check needs all seven with always(). The model job's steps run locally with this commit's build: make self-check ok (with the session's own untracked actor log set aside) and go test -trimpath ./internal/op -run TestEveryFieldInRepoLogIsKnown ok; go test -trimpath ./internal/ciworkflow/... ok. main has no branch protection (gh api 404), so ship check is the gate. The post-merge observation of the next model-only PR's CI is not yet observable.

## Links

### Out

- verifies [ci-model-only-pr](../capability/rhiza/integrity/self-check/ci-model-only-pr.md)

## History

- `01M7684F6CEBCS33NW5DDK1BVP` CI on a PR that changes only .rhiza/ (plan PRs, records, accepts, renders) runs just the model checks, not the Go suite, so a plan PR's CI takes a minute or two instead of 11-18. Owner, 2026-11-11, on the hour a milestone plan run spends after its interview: "an hour is also absurd for such a simple task", then "yes file it" for this request. Today .github/workflows/ci.yml runs lint, scripts (sharded), api (sharded), rest and budgets jobs on every PR, whatever it touches; plan PRs #435, #436 and #437 each took 15-18 minutes from interview to merge, mostly CI. Wanted: a first job decides from the PR's changed paths whether any file outside .rhiza/ changed (a paths filter, e.g. dorny/paths-filter or git diff against the merge base); when only .rhiza/ changed, CI builds rhiza and runs make self-check (rhiza check with this commit's build, including stale .rhiza/rendered) and skips the Go test jobs; any other change runs everything as today. The required status checks must still report on a model-only PR (skipped jobs count as passing, or one aggregate required job), so branch protection and rhiza ship check --settled treat it as green; check ship check and the watchdog's CI rerun still see the right checks. main's own CI (push to main) keeps running everything. Keep make check the single local entry point. Out of scope: the general CI speed-up (request 01M6Q58QG7T8V9F602YX6J2XPZ).
