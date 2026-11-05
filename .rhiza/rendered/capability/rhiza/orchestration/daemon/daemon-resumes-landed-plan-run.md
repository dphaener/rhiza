# daemon-resumes-landed-plan-run

capability · `01M6JZNFW03AN8AXVQSKVJ00RG` · parent [daemon](../daemon.md)

## acceptance

A test of daemon.Resumable (internal/daemon/pick.go) and a daemon poll testscript resume an interrupted plan run on a landed record once and journal the resume, and resume neither a plan run on a withdrawn change set nor a build run on a landed change set whose journal never reached land; the existing daemon-resumes-accepted-only, daemon-resume-once and daemon-resumes-landed-past-land tests keep passing

## statement

The daemon resumes an interrupted plan-workflow run whose change set has landed once, as it does a run on an accepted change set, while it still leaves alone a plan run whose change set is rejected, withdrawn or superseded and every other workflow's run on a change set that is not accepted, except a non-plan run on a landed change set whose journal reached its land step, which it resumes (daemon-resumes-landed-past-land)

## Links

### Out

- relates-to [plan-run-resumes-landed-record](../../workflows/built-in-workflows/plan-workflow/plan-gate/plan-run-resumes-landed-record.md) — Same exemption as rhiza run resume: plan workflow, landed state only
- decided-by [decision-plan-landed-resumable](../../../../decision/decision-plan-landed-resumable.md)

### In

- blocks [review-01m6jzdh-r3-f1](../../../../question/review-01m6jzdh-r3-f1.md)
- verifies [daemon-resumes-landed-plan-run-verified](../../../../evidence/daemon-resumes-landed-plan-run-verified.md)
- verifies [daemon-resumes-landed-plan-run-qa-1](../../../../evidence/daemon-resumes-landed-plan-run-qa-1.md)
- verifies [pr-276-fix-inbox-and-daemon-treat-an](../../../../evidence/pr-276-fix-inbox-and-daemon-treat-an.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
