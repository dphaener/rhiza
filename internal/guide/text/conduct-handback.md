# Conduct hand-back

What a conductor does when the conductor default `ship-change` hands a red step back to
it: read handback.md, spawn the role the route names, never fix anything yourself.

## handback.md

The hand-back prompt names the file. Its parts, each only when the cause has it:

- `Failing step` and `Outcome`: the step that went red and how (fail, changes, timeout).
- `## Failing commands`: each failing command; a gate's with the path of its kept output
  (`output in <path>`). Read that file in full.
- `## Output tail`: the tail of a failing command step's output, which keeps no file.
- `## Cause`: the step's own account, when it names no failing command.
- `## Failed CI runs`: each failed CI run and its saved log file.
- `## Conflicted files`: the files a conflict left.

## Routes

One per step that hands back (`handback:conduct`):

- `gate`: the failing command picks the role.
  - `rhiza run build-check`: a verifier for each leaf lacking independent Evidence; a
    builder retry when it fails.
  - `rhiza run review-check`: a critic for the check with an open blocker or major, then
    a fixer for the confirmed findings, then the checks again.
  - `rhiza check`: a fixer, its brief quoting the output as its findings.
  - `rhiza ship gate`: the gate-fix role.
- `signoff`: the owner asked for changes; a fixer with the requested changes as its
  findings, then the checks again.
- `push`: `rhiza ship open-pr`'s pre-push check failed; the gate-fix role.
- `wait-ci`: the timeout; CI never settled, so there are no failing tests or findings to
  fix and you spawn no role. Read `rhiza ship check --json`, change nothing and finish:
  the run goes through gate, signoff, push and a fresh CI wait.
- `merge`: by `## Cause` or `## Failed CI runs`.
  - A settled red CI: the ci-fix role.
  - A conflicting PR or main moved past the tested head: `rhiza ship rebase --no-push`
    (never plain `rhiza ship rebase`), then the conflict-fix role for a code conflict it
    leaves in progress.

Every route then lands and renders again on the branch (`rhiza land`, `rhiza render`,
commit what they change) and ends with `rhiza session done`.

## Limits

- `handback_cap` (3 unless the workflow sets it) caps the hand-backs: the next red after
  that parks the run stopped with an inbox item. A CI that never settles parks there.
- The steps after the merge (main-ci, close-out, ship, deploy-check) never hand back:
  nothing on the branch can fix a merged change set, so a red one stops the run.
