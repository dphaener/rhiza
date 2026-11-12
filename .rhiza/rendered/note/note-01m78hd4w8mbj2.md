# note-01m78hd4w8mbj2

note · `01M78HD4W8MBJ2RP5SFH048NMZ`

## statement

Reproduction for 'fails before the change': a delay inside the tab Update alone cannot fail the old script, since the tab Update takes the lock before the hook lands. Inject, in internal/api/session.go just before the tab Update, a wait until st.Read(sid) has Status.Hook (poll 5ms, up to 2s), and a 300ms sleep inside the Update closure before r.Tab = tab. Then the old script (HEAD~2's inbox-stalled.txtar) fails 3/3 at the first stalled check (its hook-at edit is dropped; the hook keeps the real time) and the edit-sessions one passes 3/3.

## Links

### Out

- discusses [testscripts-edit-sessions-locked](../capability/rhiza/integrity/self-check/testscripts-edit-sessions-locked.md)

## History

- `01M78ETR0NZZX29FXHDTZRYPKS` Fix flaky test scripts (4/4) / UNKNOWN STEP: failed then passed on a rerun in PR #498 at head d1e350c60c5ad977683e83207632bbe5ffa12be6; reproduce it locally, repeatedly, then fix it or quarantine it
