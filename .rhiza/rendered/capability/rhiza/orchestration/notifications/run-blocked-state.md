# run-blocked-state

capability · `01M5TYPYZ6GFFNK5BCAAHBK7K4` · parent [notifications](../notifications.md)

## acceptance

A test drives a fake session status to waiting and back and reads the run's state.

## statement

A run is blocked on the owner while its agent step's session waits for the owner's input, and stops being blocked once the owner answers.

## Links

### Out

- decided-by [status-from-hooks-first](../../../../decision/status-from-hooks-first.md)
- decided-by [m2-scope-in](../../../../decision/m2-scope-in.md)

### In

- relates-to [idle-step-blocked](../../../../question/idle-step-blocked.md) — Today only a session waiting for input (permission, question, auth) is blocked; an idle prompt is not.
- verifies [run-blocked-state-tested](../../../../evidence/run-blocked-state-tested.md)

## History

- `01M5SZJ46HYE2DDJ0STCDZE1WP` Blocked runs tell the owner: a run waiting on the owner sends a Rex notification, and genesis run status names the tab and the question.
