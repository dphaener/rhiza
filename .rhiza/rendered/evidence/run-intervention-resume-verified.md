# run-intervention-resume-verified

evidence · `01M7R3ZHACR25BS7MYQ665Q0ZH` · actor [claude-code-83c48383](../actor/claude-code-83c48383.md)

## result

PASS

## statement

Verified builder commit 4707b1d94 by re-running its Checks: block in capped scopes with RHIZA_/SHIP_/REX_ unset: go test -trimpath -count=1 -p 2 ./cmd/rhiza -run 'TestScripts/run-intervention-resume$' PASS; go test ./internal/api -run 'Intervention|CPUProfile|DaemonContinue|RunResume|RunAbandon' ok; the 52 scripts naming run resume, run abandon or intervention ok (25s); go vet ./internal/api exit 0; rhiza ship generate exit 0 (0 files changed); rhiza render exit 0 (only untracked actor records from other sessions changed index.md). The script shows rhiza run resume (queued and --foreground), resume --retry and run abandon each journalling one intervention with actor, kind, detail and step in rhiza run status --json, --dry-run of resume, retry and abandon journalling none, and RHIZA_DAEMON_RESUME (the daemon-resume path) journalling none for a resume and a retry. TestInterventionDaemonResumesJournalNothing shows daemonResume and the watchdog's queueContinue plus daemonQueuedContinues detach with RHIZA_DAEMON_RESUME=1 and journal 0. Diff read: journalling happens only in the run.resume/run.abandon op handlers, after the DryRun returns; queueContinue, the engine and run start (runEngineRetry with retry=false, never a retry) journal nothing; resume/retry/abandon behaviour is unchanged. Who journals: the TUI Inbox (run.resume, run.resume --retry, run.abandon via the op) and the web Answer action (runRunResume called directly after the last gate question is answered) both journal as person interventions; sign-off continuations go through queueContinue and journal none. Notes for the checks: the web/TUI answer-then-resume of a plan run's gate is counted as a resume intervention, arguably an owner step the run waits on; and a --foreground resume journals before runEngineRetry, so a resume that then fails still leaves its intervention.

## Links

### Out

- verifies [run-intervention-resume](../capability/rhiza/workflows/workflow-runs/run-intervention-resume.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted
