# m22-run-from-rex-ask

source · `01M6CS0C444W55N1J1W1N22HFS`

## quotes

```json
[
  {
    "by": "owner",
    "text": "Where I want to be soon: I can interact with rhiza via the mini TUI apps to examine the full health of the application, start new changes, report bugs, start new features, etc, without having this interactive session running. The daemon should be able to handle all of it. I should be in fuyll control of which parts of the process are autonomous, that should also be interactive and part of the tui.",
    "where": "orchestrator session, 2026-11-01"
  },
  {
    "by": "owner",
    "text": "make it so",
    "where": "orchestrator session, 2026-11-01"
  }
]
```

## raw_sha256

871d0cde81d89057451f1af9894c9649d362f4b34d6a800a6af8504ae9925b80

## summary

The owner's final ask for the weekend (2026-11-01 evening): run Rhiza from the TUI apps alone, with the daemon doing all the orchestration and the owner in interactive control of which parts are autonomous. Research mapped today's gap. No app starts, stops, pauses or configures the daemon; its settings are command-line flags read once (only rhiza.notifyCommand comes from git config), and it cannot be paused. The daemon starts only ship-change runs for accepted change sets; nothing starts a plan run without a human (the roadmap app's new-work keys open an interactive plan run in a tab), and the only automatic acceptance is an autonomous plan run's gate. The orchestration umbrella's accept-policy ('Acceptance policy decides which change sets are accepted without the owner') is unbuilt. There is no way to file a bug or feature from the TUI as a request, and no health screen: daemon status and rhiza status exist but say nothing of disk, the main checkout, the inbox or check together. Close-out pulls main fast-forward-only and swallows a failure, so a blocked main goes stale silently. Stalled sessions appear in the inbox but can only be opened, not nudged, and a session parked on a usage-limit prompt needs an Escape before a message reaches it. The orchestrator fixed the ship step to run make update separately (PR #184).

## title

Milestone 22: run Rhiza from Rex

## Links

### In

- derived-from [m22-run-from-rex](../decision/m22-run-from-rex.md)
- derived-from [settings-in-git-config](../decision/settings-in-git-config.md)
- derived-from [stage-autonomy-meanings](../decision/stage-autonomy-meanings.md)
- derived-from [requests-queue-local](../decision/requests-queue-local.md)
- derived-from [health-one-screen](../decision/health-one-screen.md)
- derived-from [autonomy-app-controls](../decision/autonomy-app-controls.md)

## History

- `01M6CRZP2E8H6EEEPVDAWAKQ6H` Plan milestone 22: run Rhiza from Rex, the owner's ask of 2026-11-01
