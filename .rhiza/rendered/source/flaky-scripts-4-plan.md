# flaky-scripts-4-plan

source · `01M78F2TEG8V0X15FQQX7CEMZB`

## quotes

```json
[]
```

## raw_sha256

4e8f432bbf40a71e90eb43f4c07bbe36f6e597efeae0965ddf19b3b6ec22bb2c

## session

aba0b284-b2de-4064-9020-01a1d0eda8ee

## summary

Autonomous plan run, no owner. CI run 43256685725 attempt 1 (PR #498, head d1e350c) failed scripts (4/4) in TestScripts/inbox-stalled at inbox-stalled.txtar:32: the inbox listed the session as stalled since 2026-10-24T12:00:06Z, the hook time before the script's hand edit, so the edit to hook.at (2030-01-01T00:00:00Z) was lost. Cause found by reading the code: the script rewrites .git/rhiza/sessions/<id>.json with perl and mv, without the session store's lock, as soon as the record has a hook report. The fake claude reports its hook while `rhiza session start` is still running; session start then records the tab with Store.Update, which reads the record under the flock on .<id>.lock and writes it back. When Update reads before the perl rewrite and renames after it, the rewrite is lost. inbox-permission, run-status-stalled and daemon-notify-hook hand-edit the same way. The planner asked the arbiter to choose between a locked edit command, waiting for the tab write, or quarantine; the question reached the arbiter without its options, and the arbiter ruled (flaky-scripts-4-next-step): reproduce locally with an injected delay, then fix the timing dependence or quarantine with a bug filed. The planner chose the fix: route the hand edits through a testscript command holding the store's lock.

## title

Flaky scripts (4/4) plan, autonomous run

## Links

### In

- derived-from [session-edits-under-lock](../decision/session-edits-under-lock.md)

## History

- `01M78ETR0NZZX29FXHDTZRYPKS` Fix flaky test scripts (4/4) / UNKNOWN STEP: failed then passed on a rerun in PR #498 at head d1e350c60c5ad977683e83207632bbe5ffa12be6; reproduce it locally, repeatedly, then fix it or quarantine it
