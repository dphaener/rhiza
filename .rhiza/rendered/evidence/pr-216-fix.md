# pr-216-fix

evidence · `01M6HBVFW5Q3Q9BP2Z4YAKBAWT` · actor [dphaener](../actor/dphaener.md)

## head

5b6dd154a3312c64fef40c157a9ea29290a0223f

## pr

https://github.com/dphaener/rhiza/pull/216

## result

1 check passed: check SUCCESS

## statement

Pull request #216 (https://github.com/dphaener/rhiza/pull/216) passed CI at 5b6dd15

## status

verified

## Links

### Out

- verifies [plan-run-resumes-landed-record](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/plan-run-resumes-landed-record.md)
- verifies [run-resume](../capability/rhiza/workflows/workflow-runs/run-resume.md)
- verifies [manual-resume-refuses-unaccepted](../capability/rhiza/workflows/workflow-runs/manual-resume-refuses-unaccepted.md)

## History

- `01M6GCQACNNTMW5B57HW8N7K50` Fix: rhiza run resume refuses a plan run whose record change set landed (internal/api/run.go ~792, the Terminal() guard). Every plan run lands its record during the interview, so a plan run stopped at the gate on open questions can never resume. Exempt plan-workflow runs from that guard
