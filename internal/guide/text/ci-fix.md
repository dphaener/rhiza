# CI fixer

## Task

Fix the CI findings of a ship run's red PR: `RHIZA_CHANGE` and `RHIZA_RUN` are set. You
work on the run's worktree, at the head CI failed on. Fix those findings only.

## Expected outcome

Every CI finding reproduced locally and then fixed, regenerated or quarantined, in one
new commit trailed `Review: <run>` whose `Checks:` block lists each command you ran with
its exit code. The run then pushes it and CI runs again.

## Read first

```
rhiza ship check --json
```

Its `tests` are the failing tests; its `findings` are the recognised non-test causes (a
self-check, format, vet or build-target failure), each with the failing log lines. They
name no leaf. Read the repo's contributor guide for its own scoped command for each
one: the test command for a test, the matching check command for the rest.

## Must do

Reproduce each finding locally first, with that scoped command, repeating a test until
it fails again or plainly passes. Never push to see what CI does.

- Fails locally: fix it. Stale generated files: regenerate them with the repo's own command.
- A test that does not reproduce is a flake: fix it, or quarantine it (skip it) with a
  bug filed by `rhiza request --kind bug "<flaky test and its CI log line>"`, the
  request named in a comment at the skip.

Commit only once those commands and the packages you touched pass locally. Then run
`rhiza ship generate` and `rhiza render`, fold what they change into the commit, and end
with exactly one new commit, each command's exit code in a `Checks:` block above the trailers:

```
Fix the CI findings of the export run

Checks:
- <scoped command for the finding>: exit 0
- rhiza ship generate: exit 0
- rhiza render: exit 0

Change: 01M5PWGEF85Y185JZ2ZGCNH40N
Actor: claude-code-1a2b3c4d
Review: 01M5PX3N6Y0D2W8VQK4TSE1HMB
Co-Authored-By: Agent Name <agent@example.com>
```

Keep every trailer, `Co-Authored-By:` included, in the one final paragraph, no blank line.

<!-- shared: test-scope -->

<!-- shared: unstable-ci -->

<!-- shared: command-budget nested -->

<!-- shared: outside-work nested -->

## Must not do

<!-- shared: never-list nested -->
- Push, or rerun CI, to see whether a finding is real.
- Run the repo's full suite: the PR's CI runs it again after the push.
- Touch code no finding points to.

## Context

<!-- shared: session-done -->

The ship run reaches you after a red PR whose CI named a failing test or a recognised
cause; an infra failure waits for a rerun instead. You run at most twice, so a third red
head with findings fails the run.

<!-- shared: contract-ref -->
