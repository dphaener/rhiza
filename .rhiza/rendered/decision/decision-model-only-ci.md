# decision-model-only-ci

decision · `01M768CNGSJCMZ4MAW5Y3SZ7SY`

## reasoning

Intents: `rhiza ls --kind intent` shows only milestone intents (m3 to m31: daemon, apps, Nebari, browser monitor, root, public docs, sign-off pages); none is about CI time, and the intent itself puts the general CI speed-up (request 01M6Q58) out of scope, so attaching it to a milestone would misstate the plan; a standalone request change set needs no intent. Detection: ci.yml already checks out with fetch-depth 0 and triggers only on pull_request and workflow_dispatch (no push trigger; Ship dispatches on main commits), so git diff against the merge base is enough and avoids a new action dependency, and forcing code=true on workflow_dispatch keeps main's CI full, as the intent requires. Today `check` is the aggregate branch protection and Ship read and it rejects anything but success; it must now accept skipped, but still fail when `changes` itself failed (otherwise a broken detector would skip everything and pass). A separate `model` job instead of lint avoids running self-check twice on code PRs, and matches the intent's "builds rhiza and runs make self-check". Ship: internal/api/ship.go passingConclusions already counts SKIPPED as passing in checksVerdict, and shipRerunFailed skips passing conclusions, so no code change is expected, but nothing tests it; the intent explicitly asks to check ship check and the watchdog rerun, and a leaf with an acceptance test is the cheap way to prove it. Plan review round 2 (f1, confirmed): internal/op/fields_test.go TestEveryFieldInRepoLogIsKnown is the only Go test reading the repo's own .rhiza/log, and the only other op.FieldKnown caller is internal/api/write_node.go (on write); the model job runs it so a model-only PR loses no check it has today.

## statement

(1) Change set 01M7684 serves no milestone intent: it is a standalone owner request (filed 2026-11-11, "yes file it"), like the bug and request change sets beside it. Plan its leaves under it with no serves link, and do not invent an intent or a policy.toml for it. (2) ci.yml gets a first job, `changes`, that decides with plain git (checkout already uses fetch-depth 0, no third-party action): on pull_request it lists `git diff --name-only <merge-base of base.sha and head>..<head>` and outputs code=true when any path lies outside .rhiza/, else code=false; on workflow_dispatch (Ship's runs on main commits) it always outputs code=true, so main keeps running everything. lint, scripts, api, rest and budgets run only when code=true. A new `model` job runs only when code=false: it sets up Go and runs `make self-check` (build this commit's rhiza in a temp dir, rhiza check including stale .rhiza/rendered), then `go test -trimpath ./internal/op -run TestEveryFieldInRepoLogIsKnown`: that test globs the repo's .rhiza/log and fails on any node field op.FieldKnown does not know, which fold and rhiza check never validate (only writes do), so without it a log entry written by a newer or unmerged rhiza would turn a red model-only PR green. `check` stays the single required aggregate: it needs changes, model and the five Go jobs, runs if: always(), fails unless `changes` succeeded, and passes only when every other needed job is success or skipped, with at least one of `model` or the Go jobs having succeeded. make check stays the local entry point, unchanged. (3) Yes: ship check gets its own small leaf, a Go test in internal/api, proving a check rollup with SKIPPED jobs and a SUCCESS `check` is settled and green under checksVerdict, and that shipRerunFailed reruns nothing for SKIPPED runs.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-de19e88c-1](../question/question-de19e88c-1.md)
- derived-from [model-only-ci-interview](../source/model-only-ci-interview.md)

### In

- decided-by [ci-model-only-pr](../capability/rhiza/integrity/self-check/ci-model-only-pr.md)
- decided-by [ship-check-skipped-jobs-green](../capability/rhiza/orchestration/shipping/ship-check-skipped-jobs-green.md)

## History

- `01M7684F6CEBCS33NW5DDK1BVP` CI on a PR that changes only .rhiza/ (plan PRs, records, accepts, renders) runs just the model checks, not the Go suite, so a plan PR's CI takes a minute or two instead of 11-18. Owner, 2026-11-11, on the hour a milestone plan run spends after its interview: "an hour is also absurd for such a simple task", then "yes file it" for this request. Today .github/workflows/ci.yml runs lint, scripts (sharded), api (sharded), rest and budgets jobs on every PR, whatever it touches; plan PRs #435, #436 and #437 each took 15-18 minutes from interview to merge, mostly CI. Wanted: a first job decides from the PR's changed paths whether any file outside .rhiza/ changed (a paths filter, e.g. dorny/paths-filter or git diff against the merge base); when only .rhiza/ changed, CI builds rhiza and runs make self-check (rhiza check with this commit's build, including stale .rhiza/rendered) and skips the Go test jobs; any other change runs everything as today. The required status checks must still report on a model-only PR (skipped jobs count as passing, or one aggregate required job), so branch protection and rhiza ship check --settled treat it as green; check ship check and the watchdog's CI rerun still see the right checks. main's own CI (push to main) keeps running everything. Keep make check the single local entry point. Out of scope: the general CI speed-up (request 01M6Q58QG7T8V9F602YX6J2XPZ).
