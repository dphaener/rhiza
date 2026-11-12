# question-de19e88c-1

question · `01M7689NPS6HFMGEKXZPB5MAPE`

## session

de19e88c-3c11-4405-aec4-01a1cc82df71

## statement

Which intent does change set 01M7684 (CI on a PR that changes only .rhiza/) serve? It serves none yet; there is no RHIZA_MILESTONE and no .rhiza/policy.toml. / How does ci.yml tell a model-only PR, and what runs then? / Should ship check get its own leaf proving a rollup with SKIPPED jobs and a green `check` reads as settled and green (it already treats SKIPPED as passing, untested), and that the watchdog's CI rerun ignores skipped jobs?

## Links

### In

- answers [decision-model-only-ci](../decision/decision-model-only-ci.md)

## History

- `01M7684F6CEBCS33NW5DDK1BVP` CI on a PR that changes only .rhiza/ (plan PRs, records, accepts, renders) runs just the model checks, not the Go suite, so a plan PR's CI takes a minute or two instead of 11-18. Owner, 2026-11-11, on the hour a milestone plan run spends after its interview: "an hour is also absurd for such a simple task", then "yes file it" for this request. Today .github/workflows/ci.yml runs lint, scripts (sharded), api (sharded), rest and budgets jobs on every PR, whatever it touches; plan PRs #435, #436 and #437 each took 15-18 minutes from interview to merge, mostly CI. Wanted: a first job decides from the PR's changed paths whether any file outside .rhiza/ changed (a paths filter, e.g. dorny/paths-filter or git diff against the merge base); when only .rhiza/ changed, CI builds rhiza and runs make self-check (rhiza check with this commit's build, including stale .rhiza/rendered) and skips the Go test jobs; any other change runs everything as today. The required status checks must still report on a model-only PR (skipped jobs count as passing, or one aggregate required job), so branch protection and rhiza ship check --settled treat it as green; check ship check and the watchdog's CI rerun still see the right checks. main's own CI (push to main) keeps running everything. Keep make check the single local entry point. Out of scope: the general CI speed-up (request 01M6Q58QG7T8V9F602YX6J2XPZ).
