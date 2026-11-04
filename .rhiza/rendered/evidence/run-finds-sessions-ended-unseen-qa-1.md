# run-finds-sessions-ended-unseen-qa-1

evidence · `01M6MA2PXZJZ9MGVTDWHH7P6V5` · actor [claude-code-3a61d2cb](../actor/claude-code-3a61d2cb.md)

## result

PASS

## statement

QA with a binary built from the branch head (go build ./cmd/rhiza) in scratch repos, a fake claude, and a rex wrapper that sends block inspect to the real Rex on this VM. Started an agent-step run, killed it after its session was recorded, then ran rhiza run resume --foreground. (1) Real Rex answered 'block target ... not found' for the session's block: the record got exit {state: lost, detail: 'Rex block not found'}, the step failed, rhiza run status showed 'lost: Rex block not found', and rex block close was called for the pane. (2) Block present, then a SessionEnd hook (rhiza session hook SessionEnd) delivered 2s into the resume: the next poll wrote a lost exit 'SessionEnd without an exit report' with no block inspect call after the rejoin's, and the step failed. (3) Review run (run-resume-lost-fixer fixtures) killed while its change-set fixer was in flight, resumed with real Rex answering not found: fixer attempt 1 judged 'fail (lost: Rex block not found)', attempt 2 started and committed, check alpha re-ran in cycle 2, run passed (1 of 1 checks).

## Links

### Out

- verifies [run-finds-sessions-ended-unseen](../capability/rhiza/workflows/workflow-runs/run-finds-sessions-ended-unseen.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
