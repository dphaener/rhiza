# plan-resume-landed-record-milestone-6

decision · `01M6GCZW3NH9ZBYN4389ZGDWCB`

## reasoning

Intent: plan-run-stops-open's change set served both milestone-4 and milestone-6; milestone-6 is the open milestone about the daemon running and planning its own PRs unattended, and a plan run that can never resume after the gate stops it is exactly what breaks unattended operation, so milestone-6 is the better single home. Narrowness: internal/api/run.go guards with changeset.State.Terminal(), which covers landed, rejected, withdrawn and superseded. Only landing is expected for a plan run (its record lands during the interview); a rejected, withdrawn or superseded change set still means the owner or the model abandoned the work, and resuming then would ship a dead plan. Keying on the workflow name, which run start --new already checks via planWorkflow, keeps build and review runs protected by the guard as before.

## statement

Change set 01M6GCQA serves milestone-6 (Rhiza runs unattended overnight). The exemption in run resume's change-set guard is narrow: a run of the plan workflow (def.Name == planWorkflow) whose change set is landed is let through; every other terminal state (rejected, withdrawn, superseded) still refuses a plan run, and runs of every other workflow keep the guard unchanged

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-b693317e-1](../question/question-b693317e-1.md)
- derived-from [plan-resume-landed-record-interview](../source/plan-resume-landed-record-interview.md)

### In

- decided-by [plan-run-resumes-landed-record](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/plan-run-resumes-landed-record.md)

## History

- `01M6GCQACNNTMW5B57HW8N7K50` Fix: rhiza run resume refuses a plan run whose record change set landed (internal/api/run.go ~792, the Terminal() guard). Every plan run lands its record during the interview, so a plan run stopped at the gate on open questions can never resume. Exempt plan-workflow runs from that guard
