# tui-small-apps-interview

source · `01M6MYKNE2EZ4DA6RS0B5XEGHQ`

## quotes

```json
[
  {
    "by": "owner",
    "text": "All of the apps are pretty bad and need an overhaul, and most of them are pretty useless. We're going to have a full milestone for that but not until we get everything running on the VM.",
    "where": "2026-11-03, from the orchestrator's notes"
  },
  {
    "by": "owner",
    "text": "1. there are too many. i won't use most of them\n2. they aren't using any of the form elements from rexuikit. they don't look great\n3. i want to refine what we do keep to ensure i can fully manage rhiza runs with as little tui surface as possible.",
    "where": "2026-11-05, typed answer 1: what is wrong with the apps"
  },
  {
    "by": "owner",
    "text": "multiple small apps. single purpose. i don't want the tui to dominate the terminal. think in and out. the proposed actions seem good.",
    "where": "2026-11-05, typed answer 2: the shape of the apps and the proposed actions"
  },
  {
    "by": "agent (picked by the owner)",
    "text": "Delete them: remove the code and Rex actions of every app not in the new set; the CLI keeps every command",
    "where": "2026-11-05, pick: old apps"
  },
  {
    "by": "agent (picked by the owner)",
    "text": "Five apps: Inbox (pick an item and do its one action: answer a question or gate, sign off a run, retry or abandon a failed run, then close); Request (file a feature, bug or milestone); Runs (a read-only glance at live runs and their current steps); Daemon (pause or resume, and the run limit); Autonomy (the plan, accept, ship and sign-off toggles)",
    "where": "2026-11-05, pick: app set"
  },
  {
    "by": "agent (picked by the owner)",
    "text": "Small modal: each app opens as a rexuikit Modal sized to its content (Form, List or Dialog) over the current pane, and closes on Esc or when its job is done; never full-screen",
    "where": "2026-11-05, pick: look"
  },
  {
    "by": "agent (picked by the owner)",
    "text": "TUI first: the TUI milestone is 26, Nebari moves to 27",
    "where": "2026-11-05, pick: order"
  },
  {
    "by": "agent (picked by the owner)",
    "text": "Plan all 7 change sets now",
    "where": "2026-11-05, pick: depth"
  }
]
```

## raw_sha256

7964ae58c8f339c243ff0f043a021fd2f067e95ecb7ea881dcd7548b7b4f080c

## summary

Interview on 2026-11-05, in chat with the orchestrator, planning the TUI overhaul the owner had deferred on 2026-11-03 until Rhiza ran on the VM. The owner said the ten Rex apps are too many, do not use rexuikit's form elements and look poor, and asked for the kept apps to be refined so Rhiza runs can be fully managed from as little TUI surface as possible. The owner chose multiple small single-purpose apps that do not dominate the terminal ('think in and out') and approved the proposed actions for managing runs from the TUI: act on the inbox (answer questions and gates, sign off, retry or abandon failed runs), watch live runs, file a request, and steer the daemon (pause, run limit, autonomy toggles); the orchestrator had done all of these by hand that night. From the agent's options the owner picked: delete every app not in the new set (the CLI keeps every command); a five-app set (Inbox, Request, Runs, Daemon, Autonomy); each app a small rexuikit Modal sized to its content that closes on Esc or when its job is done, never full-screen; the TUI milestone first, as milestone 26, with the Nebari import moved to 27; and all seven change sets planned now. This reverses the earlier decision tui-many-mini-apps. Left open: which inbox kinds the Inbox app acts on beyond those named, whether a sign-off in it includes requesting changes, and how the Rex action opens the modal.

## title

Milestone 26: small single-purpose TUI apps, with the owner

## Links

### In

- derived-from [tui-single-purpose-apps](../decision/tui-single-purpose-apps.md)
- derived-from [tui-run-actions](../decision/tui-run-actions.md)
- derived-from [tui-overhaul-why](../decision/tui-overhaul-why.md)
- derived-from [tui-rexuikit-components](../decision/tui-rexuikit-components.md)
- derived-from [tui-five-apps](../decision/tui-five-apps.md)
- derived-from [tui-delete-old-apps](../decision/tui-delete-old-apps.md)
- derived-from [tui-small-modal](../decision/tui-small-modal.md)
- derived-from [tui-milestone-first](../decision/tui-milestone-first.md)
- derived-from [tui-plan-all-now](../decision/tui-plan-all-now.md)
- derived-from [tui-modal-go-opener](../decision/tui-modal-go-opener.md)

## History

- `01M6MYK80HN9E4DG8QREZH38DW` Milestone 26: Rhiza runs are fully managed from five small single-purpose Rex apps built on rexuikit, each a small modal that opens, does one job and closes; the Nebari import moves to milestone 27
