# session-edits-under-lock

decision · `01M78F38R9PH6BXX9STFQ3TKSM`

## reasoning

inbox-stalled lost its hand edit of hook.at when session start's locked tab Update read the record before the unlocked rewrite and wrote it back after. Taking the same lock serialises the edit with every locked writer, now and later; waiting for the tab write only covers today's writers, and quarantine leaves the four scripts' checks unrun. This carries out the arbiter's ruling flaky-scripts-4-next-step (fix the timing dependence).

## statement

Testscripts that hand-edit a live run's session record make the edit through a testscript command that holds the session store's flock on .<id>.lock, the lock Store.Update takes, instead of an unlocked perl and mv; no script is quarantined

## Links

### Out

- derived-from [flaky-scripts-4-plan](../source/flaky-scripts-4-plan.md)

### In

- decided-by [testscripts-edit-sessions-locked](../capability/rhiza/integrity/self-check/testscripts-edit-sessions-locked.md)

## History

- `01M78ETR0NZZX29FXHDTZRYPKS` Fix flaky test scripts (4/4) / UNKNOWN STEP: failed then passed on a rerun in PR #498 at head d1e350c60c5ad977683e83207632bbe5ffa12be6; reproduce it locally, repeatedly, then fix it or quarantine it
