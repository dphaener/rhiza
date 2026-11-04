# lost-session-judged-qa-1

evidence · `01M6MA2TNK59Q1RGGD168037KQ` · actor [claude-code-3a61d2cb](../actor/claude-code-3a61d2cb.md)

## result

PASS

## statement

QA with the built binary, resuming a killed run whose session had a SessionEnd hook mark and no exit. Plain agent step: step failed, detail 'lost: SessionEnd without an exit report'. Step with until = 'test -f ok', until_when = 'done', done mark set via rhiza session done: with ok in the run worktree the step passed ('lost: SessionEnd without an exit report and until held', resume exit 0); without it the step failed ('... but until did not hold'). In every case the record's exit has state lost and no code, and rex block close was called for the session's pane.

## Links

### Out

- verifies [lost-session-judged](../capability/rhiza/workflows/workflow-runs/lost-session-judged.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
