# Review

A review run checks a built change set: fresh check sessions, a critic per check, fixers.
This topic is an index of eleven role topics. Load yours, as a skill or with
`rhiza guide <topic>`:

A conductor run's three checks are check-plan, check-code and check-qa. The classic workflows (ship-change, review, docs-review) keep the older four.

- Plan check (conductor): `/rhiza-check-plan` (`rhiza guide check-plan`). A check named `scope`
  loads this one too.
- Code check (conductor): `/rhiza-check-code` (`rhiza guide check-code`).
- Compliance check (classic): `/rhiza-check-compliance` (`rhiza guide check-compliance`). A check
  named `scope` loads this one too.
- Quality check (classic): `/rhiza-check-quality` (`rhiza guide check-quality`).
- QA check: `/rhiza-check-qa` (`rhiza guide check-qa`).
- Model check (classic): `/rhiza-check-model` (`rhiza guide check-model`).
- Critic: `/rhiza-critic` (`rhiza guide critic`).
- Fixer: `/rhiza-fixer` (`rhiza guide fixer`).
- CI fixer of a ship run: `/rhiza-ci-fix` (`rhiza guide ci-fix`).
- Gate fixer of a ship run: `/rhiza-gate-fix` (`rhiza guide gate-fix`).
- Conflict fixer of a ship run: `/rhiza-conflict-fix` (`rhiza guide conflict-fix`).

Each role ends with `rhiza session done` as its own command, as its topic says.
