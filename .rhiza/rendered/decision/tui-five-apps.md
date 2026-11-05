# tui-five-apps

decision · `01M6MZ1M8T2A92HB1Q7BGG19MM`

## statement

The TUI is five apps: Inbox (pick an item and do its one action: answer a question or gate, sign off a run, retry or abandon a failed run, then close), Request (file a feature, bug or milestone), Runs (a read-only glance at live runs and their current steps), Daemon (pause or resume, and the run limit) and Autonomy (the plan, accept, ship and sign-off toggles).

## wording

Inbox: pick an item and do its one action (answer a question or gate, sign off a run, retry or abandon a failed run), then close. Request: file a feature, bug or milestone. Runs: a read-only glance at live runs and their current steps. Daemon: pause or resume, and the run limit. Autonomy: the plan, accept, ship and sign-off toggles, the agent's option, picked by the owner

## Links

### Out

- derived-from [tui-small-apps-interview](../source/tui-small-apps-interview.md)

### In

- decided-by [m26-inbox-app](../intent/milestone-26/m26-inbox-app.md)
- decided-by [m26-request-app](../intent/milestone-26/m26-request-app.md)
- decided-by [m26-runs-app](../intent/milestone-26/m26-runs-app.md)
- decided-by [m26-daemon-app](../intent/milestone-26/m26-daemon-app.md)
- decided-by [m26-autonomy-app](../intent/milestone-26/m26-autonomy-app.md)
- decided-by [inbox-app-pick-list](../capability/rhiza/inbox-app-pick-list.md)
- decided-by [inbox-app-answer](../capability/rhiza/inbox-app-answer.md)
- decided-by [inbox-app-signoff](../capability/rhiza/inbox-app-signoff.md)
- decided-by [inbox-app-retry-abandon](../capability/rhiza/inbox-app-retry-abandon.md)
- decided-by [request-app-form](../capability/rhiza/request-app-form.md)
- decided-by [runs-app-glance](../capability/rhiza/runs-app-glance.md)
- decided-by [daemon-app-pause](../capability/rhiza/daemon-app-pause.md)
- decided-by [daemon-app-run-limit](../capability/rhiza/daemon-app-run-limit.md)
- decided-by [autonomy-app-toggles](../capability/rhiza/autonomy-app-toggles.md)
- decided-by [rex-install-five-apps](../capability/rhiza/rex-install-five-apps.md)

## History

- `01M6MYK80HN9E4DG8QREZH38DW` Milestone 26: Rhiza runs are fully managed from five small single-purpose Rex apps built on rexuikit, each a small modal that opens, does one job and closes; the Nebari import moves to milestone 27
