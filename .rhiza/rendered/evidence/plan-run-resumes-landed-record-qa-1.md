# plan-run-resumes-landed-record-qa-1

evidence · `01M6H6HW88ZY4X60T96Q58NCSX` · actor [claude-code-1c059d7a](../actor/claude-code-1c059d7a.md)

## result

PASS

## statement

QA by hand with this branch built to a scratch binary, in a fresh git repo with RHIZA_* env cleared: a plan workflow (one wait step on ../go) started with run start --foreground --no-wait stopped at step one; rhiza land landed its record change set (state landed); after touching ../go, rhiza run resume --foreground <run> continued visit 1, passed, and ended succeeded with exit 0 (run status --json state succeeded). In the same repo a hand-written ship-change journal on the landed change set was refused with exit 3, "is landed, no longer accepted" and next: rhiza run abandon BUILDLAND01, journal byte-identical; a plan journal on a rejected change set was refused with exit 3, "is rejected, no longer accepted", next: rhiza run abandon PLANREJ01, journal unchanged. As a control, the merge-base binary run through the same script refused the landed plan run (exit 3, run stayed stopped).

## Links

### Out

- verifies [plan-run-resumes-landed-record](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/plan-run-resumes-landed-record.md)

## History

- `01M6GCQACNNTMW5B57HW8N7K50` Fix: rhiza run resume refuses a plan run whose record change set landed (internal/api/run.go ~792, the Terminal() guard). Every plan run lands its record during the interview, so a plan run stopped at the gate on open questions can never resume. Exempt plan-workflow runs from that guard
