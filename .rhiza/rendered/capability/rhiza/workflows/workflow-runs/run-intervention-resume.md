# run-intervention-resume

capability · `01M7R2FY1KB5VKZZYS1XGDHRBC` · parent [workflow-runs](../workflow-runs.md)

## acceptance

Testscripts cmd/rhiza/testdata/script/run-intervention-resume.txtar: rhiza run resume, rhiza run resume --retry and rhiza run abandon on runs each show the one matching intervention in rhiza run status --json; (failure case) a run resumed by the daemon (daemon-resume path) and one retried by the watchdog show intervention_count 0, covered by go test ./internal/api -run 'Intervention' or the script

## files

internal/api/run.go, internal/api/daemon_watchdog.go, internal/api/run_intervention_test.go, cmd/rhiza/testdata/script/run-intervention-resume.txtar

## must_not

Change what resume, retry or abandon do; count a --dry-run

## statement

A person's rhiza run resume journals a resume intervention, with --retry (or rhiza run start retrying a failed run of the change set) a retry intervention, and rhiza run abandon an abandon intervention, each with the actor; the daemon's own resume and the watchdog's retry journal none

## Links

### Out

- depends-on [run-intervention-journal](run-intervention-journal.md)
- decided-by [intervention-is-a-person-command](../../../../decision/intervention-is-a-person-command.md)

### In

- verifies [run-intervention-resume-verified](../../../../evidence/run-intervention-resume-verified.md)
- verifies [run-intervention-resume-qa-1](../../../../evidence/run-intervention-resume-qa-1.md)
- verifies [run-intervention-resume-qa-2](../../../../evidence/run-intervention-resume-qa-2.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted
