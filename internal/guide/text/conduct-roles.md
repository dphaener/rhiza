# Conduct roles

Every role a conductor spawns, the flags each takes and how check roles work.

## Spawning

```
rhiza run spawn --role <role> --brief <file> [--leaf <leaf>] [--check <check>] [--model <model>] [--no-wait]
```

- `--brief <file>` always: the child's brief.md holds the contract section, then this file.
- `--leaf <leaf>` for a builder and a verifier: the leaf, by handle.
- `--check <check>` for a critic: the check whose latest round it judges.
- `--model <model>` optional: the model the child's agent starts on; default the agent's.
- `--no-wait` for a conductor: return once the child is started, with its session and
  no outcome.

Without `--no-wait` the spawn waits for the child (held while a usage or machine hold
is on) and prints `pass` when it said done or exited 0, else `fail`; that default is for
other callers. Under `--no-wait` the run's engine records the child's end, closes its
pane if it said done and wakes the parent session with a message: `Your spawned child
session <session> ended.`, then its role, leaf, outcome and detail and a pointer to
`rhiza run status <run> --json`. A conductor ends its turn after spawning and acts on
each such message. Only a builder's and a verifier's brief.md carry the run's notes.

## Roles

Each child follows `rhiza guide <role>` and the skill of the same name, `/rhiza-<role>`.

- `builder` (`rhiza guide builder`): one per leaf, with `--leaf`; again as a retry when the
  leaf fails verification, its brief quoting the verifier's findings verbatim.
- `verifier` (`rhiza guide verifier`): after each builder, with `--leaf`.
- `check-plan` (`rhiza guide check-plan`): once every leaf is verified, and after a fixer.
- `check-code` (`rhiza guide check-code`): beside check-plan.
- `check-qa` (`rhiza guide check-qa`): beside check-plan.
- `critic` (`rhiza guide critic`): with `--check <check>`, only for a check whose latest
  round has a blocker or major finding.
- `fixer` (`rhiza guide fixer`): for the critic's confirmed findings, then the checks again;
  on a hand-back, for `rhiza check` output or the owner's requested changes.
- `gate-fix` (`rhiza guide gate-fix`): on a hand-back from a failed `rhiza ship gate`, the
  pre-push check of `rhiza ship open-pr`, or the local gate or pre-push check failing on
  a merged head (pr-green, rebase).
- `ci-fix` (`rhiza guide ci-fix`): on a hand-back from a settled red CI at pr-green.
- `conflict-fix` (`rhiza guide conflict-fix`): for a code conflict the ship steps leave in
  progress when they merge origin/main into the branch (the CI wait, the rebase step).

## Check roles

- An agent step that declares `[[step.<name>.check]]` tables is a conductor step: its
  session and every child it spawns put their questions to the arbiter, whether or not
  the run is autonomous.
- A `check-<check>` role sets the child's `RHIZA_CHECK` to `<check>` and `RHIZA_CYCLE` to
  the step's latest cycle, else the round it records. A critic gets both from `--check`,
  its cycle the round it judges.
- `--check` naming a check the step does not declare is refused with the declared names
  (`name one of its checks: plan, code, qa`); outside a step declaring checks, any
  `--check` is refused.
- The conductor default `ship-change` declares `plan`, `code` and `qa` on its conduct step.
  `compliance`, `quality` and `model` belong to the classic workflows (`review`,
  `ship-change-classic`), whose checks step runs them itself: a conductor never spawns them.
