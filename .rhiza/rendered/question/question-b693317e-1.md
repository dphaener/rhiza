# question-b693317e-1

question · `01M6GCWANVBSDT4Y051JSJJWXP`

## session

b693317e-dc4d-40eb-b2ec-01a1a0cbe9f7

## statement

Change set 01M6GCQA (Fix: run resume refuses a plan run whose record landed) serves no intent yet. Which landed intent should it serve? The broken capability is plan-run-stops-open (a plan run stopped on open questions resumes once they are answered), whose change set served milestone-4 and milestone-6. / How narrow should the exemption in run resume's change-set guard (internal/api/run.go ~791) be for plan runs?

## Links

### In

- answers [plan-resume-landed-record-milestone-6](../decision/plan-resume-landed-record-milestone-6.md)

## History

- `01M6GCQACNNTMW5B57HW8N7K50` Fix: rhiza run resume refuses a plan run whose record change set landed (internal/api/run.go ~792, the Terminal() guard). Every plan run lands its record during the interview, so a plan run stopped at the gate on open questions can never resume. Exempt plan-workflow runs from that guard
