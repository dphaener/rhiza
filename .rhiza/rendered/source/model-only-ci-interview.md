# model-only-ci-interview

source · `01M768DNBE2DQ87QQQJMK969AT`

## quotes

```json
[]
```

## raw_sha256

dac111b27f90e3325faf5e4416d7dfd16e0f0213bad7e9be0cdcd6f5bb17d99b

## session

de19e88c-3c11-4405-aec4-01a1cc82df71

## summary

Autonomous plan run for change set 01M7684 (CI on a PR that changes only .rhiza/). The planner read .github/workflows/ci.yml (lint, scripts x4, api x2, rest, budgets and an aggregate check job; triggers pull_request and workflow_dispatch on a sha, checkout fetch-depth 0), the Makefile targets (lint = fmt-check + vet + self-check), internal/ciworkflow's tests of ci.yml, and ship check's checksVerdict (SKIPPED already counts as passing, untested) and shipRerunFailed. It raised three choices as one question: the intent to serve, how ci.yml detects a model-only PR and what runs then, and whether ship check gets its own test leaf. No owner took part; the arbiter answered with Decision decision-model-only-ci: no serves link (a standalone owner request); a plain-git `changes` job (code=true on any path outside .rhiza/, always true on workflow_dispatch) gating the five Go jobs, a `model` job running make self-check when code=false, and `check` as the single aggregate passing on success-or-skipped with changes succeeded; and a separate internal/api test leaf for the SKIPPED verdict and rerun.

## title

Plan interview: CI on a model-only PR (autonomous)

## Links

### In

- derived-from [decision-model-only-ci](../decision/decision-model-only-ci.md)

## History

- `01M7684F6CEBCS33NW5DDK1BVP` CI on a PR that changes only .rhiza/ (plan PRs, records, accepts, renders) runs just the model checks, not the Go suite, so a plan PR's CI takes a minute or two instead of 11-18. Owner, 2026-11-11, on the hour a milestone plan run spends after its interview: "an hour is also absurd for such a simple task", then "yes file it" for this request. Today .github/workflows/ci.yml runs lint, scripts (sharded), api (sharded), rest and budgets jobs on every PR, whatever it touches; plan PRs #435, #436 and #437 each took 15-18 minutes from interview to merge, mostly CI. Wanted: a first job decides from the PR's changed paths whether any file outside .rhiza/ changed (a paths filter, e.g. dorny/paths-filter or git diff against the merge base); when only .rhiza/ changed, CI builds rhiza and runs make self-check (rhiza check with this commit's build, including stale .rhiza/rendered) and skips the Go test jobs; any other change runs everything as today. The required status checks must still report on a model-only PR (skipped jobs count as passing, or one aggregate required job), so branch protection and rhiza ship check --settled treat it as green; check ship check and the watchdog's CI rerun still see the right checks. main's own CI (push to main) keeps running everything. Keep make check the single local entry point. Out of scope: the general CI speed-up (request 01M6Q58QG7T8V9F602YX6J2XPZ).
