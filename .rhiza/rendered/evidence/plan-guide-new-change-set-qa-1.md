# plan-guide-new-change-set-qa-1

evidence · `01M6A8S58AQ3AYTGEW6DG0BF3G` · actor [claude-code-b0fb6bcc](../actor/claude-code-b0fb6bcc.md)

## result

PASS

## statement

QA at cafec91c: `rhiza guide plan` from the scratch build prints the new step-2 text for an intention change set that serves no intent yet. Followed it by hand on the draft that run start --new opened: change show --json gave "serves": []; added an intent with rhiza apply, then ran the guide's exact `rhiza link <cs> serves <intent> --change <cs>`; change show --json then listed that intent's ID under serves.

## Links

### Out

- verifies [plan-guide-new-change-set](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-guide-new-change-set.md)

## History

- `01M686RVAD1EHGERBBDJN3S45G` The roadmap app starts a plan run for a new milestone or a new change set from a typed one-line intent
