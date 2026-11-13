# Gate fixer

## Task

Fix the failed local gate or pre-push check of a ship run: `RHIZA_CHANGE` and `RHIZA_RUN`
are set. You work on the run's worktree, at the head the gate failed on. CI is unstable,
so the project's local gate (`rhiza.ship.localGate`) runs before the push; the pre-push
check (`rhiza check`, then `rhiza.ship.prePushCheck`) runs on the head before a push.
Fix what one of them reported, and nothing else.

## Expected outcome

The failure reproduced locally and then fixed or quarantined, in one new commit trailed
`Review: <run>` whose `Checks:` block lists each command you ran with its exit code. The
run then runs the gate again.

## Read first

```
rhiza run status --json
```

Read the failed gate step's kept output and, for a pre-push failure, its failing lines.
The kept output ends with a `Full log:` line naming the whole gate log file: read that
file for the full failure. Read the repo's contributor guide for the scoped command of
each failing test.

## Must do

Never rerun `rhiza ship gate`: it reruns the whole gate. Reproduce the failure with the
failing test's own scoped command, repeating it until you see it
fail again or plainly pass.

- Fails locally: fix it.
- A test that does not reproduce is a flake: fix it, or quarantine it (skip it) with a
  bug filed by `rhiza request --kind bug "<flaky test and its failure line>"`, the
  request named in a comment at the skip.

Commit only once the failing tests and the packages you touched pass locally. Then run
`rhiza ship generate` and `rhiza render`, fold what they change into the commit, and end
with exactly one new commit, each command's exit code in a `Checks:` block above the trailers:

```
Fix the failed gate of the export run

Checks:
- <scoped command for the failing test>: exit 0
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
- Rerun `rhiza ship gate`, or push, to see whether the failure is real.
- Run the repo's full suite: the gate runs it again.
- Touch code the failure does not point to.

## Context

<!-- shared: session-done -->

The ship run reaches you after its local gate or pre-push check failed before a push. You
run at most twice, so a third failed head fails the run.

<!-- shared: contract-ref -->
