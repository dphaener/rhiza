# Conduct

## Task

Conduct one change set in one session: `RHIZA_CHANGE` and `RHIZA_RUN` are set. You plan
nothing yourself beyond missing leaves; you spawn the sessions that build, verify, check,
judge and fix, and you land the result on the branch. Each phase is a command; brief
each child tightly (the leaf, the one thing to do, nothing else).

## Expected outcome

Every leaf built and verified, the plan and code and QA checks passed, the change set
landed and rendered on the branch, then `rhiza session done`.

## Read first

```
rhiza slice --change "$RHIZA_CHANGE" --json
rhiza run status --json
```

Read the slice, then the code it touches. A change set may be part-way through: read
what is done and carry on from there.

## Must do

Run these phases in order.

- Plan: when the slice has no leaves, write them with `rhiza apply` under the autonomous
  rules: no owner is there to ask, so a question goes to the arbiter and its Decision
  stands. Load `rhiza guide planner`.
- Build: for each leaf, `rhiza run spawn --role builder --leaf <leaf> --brief <file>`,
  then `rhiza run spawn --role verifier --leaf <leaf> --brief <file>`. A leaf that fails
  verification gets a builder retry with the findings in its brief.
- Check: spawn `--role check-plan`, `--role check-code` and `--role check-qa`, each with
  `--brief <file>`; each loads its own topic and records a round with `rhiza run findings --check`.
- Critic: spawn `--role critic` only for a check whose latest round has a blocker or
  major finding; a clean or minor-only round needs none.
- Fix: spawn `--role fixer` for the confirmed findings, then run the checks again.
- Gate: `rhiza run review-check` must pass before you land; it is the run's gate, never yours to waive.
- Land: on the branch, `rhiza land`, then `rhiza render`, and commit what they change.
- Finish with `rhiza session done`.

<!-- shared: command-budget nested -->

## Must not do

<!-- shared: never-list nested -->
- Merge or push: the ship run does both.
- Edit `.rhiza/log`: only the CLI writes it.
- Pass your own gate: a check's round and verdict come from its own session.
- Build, check or fix yourself: spawn the role.

## Context

<!-- shared: session-done -->

<!-- shared: contract-ref -->
