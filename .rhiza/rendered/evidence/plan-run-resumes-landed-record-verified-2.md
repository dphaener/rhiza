# plan-run-resumes-landed-record-verified-2

evidence · `01M6H9X1JYT00GBVBEAG1BDMH4` · actor [claude-code-7b802f75](../actor/claude-code-7b802f75.md)

## result

PASS

## statement

At 34815079 (attempt 2, fix for model finding f1): read the fix diff; it only amends the to-be statements of run-resume, manual-resume-refuses-unaccepted and Decision manual-resume-refuses to name the exception (a plan-workflow run whose change set has landed is resumed), keeping the original wording otherwise; rhiza show --change confirms all three now read that way, so the to-be model no longer states two contradictory rules. Code is unchanged since 92222406 (internal/api/run.go exempts only landed and only for workflow plan). go test ./cmd/rhiza -count=1 -run 'TestScripts/(run-resume-landed-plan|run-resume-unaccepted|run-resume)$' passes: the plan run resumes past its wait after its record landed and succeeds; the landed ship-change run and the withdrawn plan run each exit non-zero naming the state and 'next: rhiza run abandon <id>', journals cmp-identical. fmt-check, go vet and make self-check pass. Full go test ./... failed only on TestReadBudget, TestBudget (timing) and arbiter-agent-tabs under load average ~13; arbiter-agent-tabs passes rerun alone, and this commit touches only .rhiza/ files.

## Links

### Out

- verifies [plan-run-resumes-landed-record](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/plan-run-resumes-landed-record.md)

## History

- `01M6GCQACNNTMW5B57HW8N7K50` Fix: rhiza run resume refuses a plan run whose record change set landed (internal/api/run.go ~792, the Terminal() guard). Every plan run lands its record during the interview, so a plan run stopped at the gate on open questions can never resume. Exempt plan-workflow runs from that guard
