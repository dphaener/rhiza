# autonomous-review

capability · `01M5WRHEQX9FSJ95GVE3BH7EQJ` · parent [autonomous-runs](../autonomous-runs.md)

## acceptance

A testscript runs an autonomous workflow on a fake adapter that accepts a change set under rule autonomous and raises one question its fake arbiter answers: run status --json lists the accept (change set and entry) and the question with its Decision's handle, and the markdown shows both; a normal run's status lists neither.

## statement

rhiza run status on an autonomous run lists what the run did in the owner's place: each accept it made under rule autonomous, and each question its arbiter decided with the Decision that answers it.

## Links

### Out

- depends-on [autonomous-gate-rule](autonomous-gate-rule.md)
- depends-on [arbiter-on-question](../arbiter-sessions/arbiter-on-question.md)
- decided-by [autonomous-mode](../../../../decision/autonomous-mode.md)
- decided-by [autonomous-accepts-journalled](../../../../decision/autonomous-accepts-journalled.md)

### In

- relates-to [autonomous-undo](../../../../question/autonomous-undo.md) — The autonomous-mode Decision says all of it is reviewable and undoable; this change set makes it reviewable, and how it is undone is open.
- verifies [autonomous-review-verified](../../../../evidence/autonomous-review-verified.md)
- verifies [pr-39-autonomous-mode-and-the](../../../../evidence/pr-39-autonomous-mode-and-the.md)

## History

- `01M5WPTVT93DTC076WJNB7E2KX` Autonomous mode and the arbiter: a run can skip the owner's gates and questions, recording them as rule autonomous and (assumed, would have escalated), and the engine calls a fresh arbiter session at a decision point; agent permissions are never bypassed.
