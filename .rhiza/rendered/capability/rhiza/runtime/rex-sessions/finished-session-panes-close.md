# finished-session-panes-close

capability · `01M6H3V1Z1MTKF0685QFDQJZAZ` · parent [rex-sessions](../rex-sessions.md)

## acceptance

Testscripts with the fake rex show a rex block close of each finished session's block at the moment its step judges it, before the next session or step starts, not only at close-out: for a review checks step (checks, critic, fixer, verifier), a build step (failed and passed builder, verifier), and an agent step that said done with until held (plan steps); a Go test shows a failed arbiter that said done has its pane closed while a blocked one keeps it

## statement

Every run session's pane closes as soon as its driver records the session finished, whatever the step kind (leaf builder or verifier, review check, critic, change-set or leaf fixer, plan interview, review, critic or reviser, arbiter), so a run's tabs hold only sessions still working or waiting on someone, and a tab goes once its last pane closes

## Links

### Out

- decided-by [finished-panes-close-at-once](../../../../decision/finished-panes-close-at-once.md)

## History

- `01M6H3TC786A6WVKRXC0WQ1N3K` Agent panes close when their session finishes: every driver closes a finished session's pane at once (review checks, critic, fixers, plan steps), not at close-out
