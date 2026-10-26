# agent-step-done-signal

capability · `01M5X3GC1F8YT2560E864QREFK` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript with a fake agent that goes idle, opens a PR and only later runs rhiza session done: the build step does not pass while idle, and passes after the done mark.

## statement

An agent step's until can wait for the agent to say it is done: rhiza session done marks the session's record, and ship-change's build passes only on that mark (or exit 0) and an open PR, never on idle alone. Its prompt tells the agent to run rhiza session done as its very last action, after the owner's sign-off.

## Links

### Out

- decided-by [ship-change-holds-up](../../../../decision/ship-change-holds-up.md)

### In

- verifies [agent-step-done-signal-tests](../../../../evidence/agent-step-done-signal-tests.md)

## History

- `01M5X3GBH4ZQTZVXFEHR966747` ship-change that holds up: the agent step passes only when the agent says it is done, land renders with the branch's own build in Rhiza's repo, a PR that main has moved past is rebased, and a failed run can be resumed from the step that failed.
