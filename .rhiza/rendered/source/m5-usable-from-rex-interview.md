# m5-usable-from-rex-interview

source · `01M6731MPG929PA11VHEV8D5R0`

## quotes

```json
[
  {
    "by": "owner",
    "text": "What I'd like to do is find the fastest possible route to this being usable, by me, from rex, using the TUI to launch things.",
    "where": "orchestrator session, 2026-10-30"
  },
  {
    "by": "owner",
    "text": "Currently the commands and mini apps don't provide a lot of value. Forinstance, the roadmap. It just prints a json blob of things and isn't actionable in any way.",
    "where": "orchestrator session, 2026-10-30"
  }
]
```

## raw_sha256

f54a3bb1ebb73725e311ae32293e796c2843a0b908d59a4558e4ebcb3e4f31c8

## summary

After milestone 4 landed, the owner asked for the fastest route to Rhiza being usable by them from Rex, launching things from the TUI. The mini apps give little value today: each read app prints its operation's JSON in a scrolling pane, and nothing in them is actionable. The orchestrator proposed making two apps the hub: the roadmap app as a cursor list of milestones and change sets whose selected change set can be opened, planned, shipped, accepted, or its blocking questions answered; and a new inbox app over the existing inbox operation, where questions, sign-offs and failed runs are acted on in place. The other mini apps stay as they are, keeping the many-mini-apps decision. The owner picked 'Roadmap + Inbox' as the first slice and chose to build it through the daemon's ship runs rather than by the orchestrator directly.

## title

Milestone 5: Rhiza usable from Rex

## Links

### In

- derived-from [m5-roadmap-inbox-first](../decision/m5-roadmap-inbox-first.md)
- derived-from [m5-built-by-daemon](../decision/m5-built-by-daemon.md)
- derived-from [m5-other-apps-stay](../decision/m5-other-apps-stay.md)
- derived-from [app-runs-open-rex-tabs](../decision/app-runs-open-rex-tabs.md)
- derived-from [app-plan-interactive-ship-autonomous](../decision/app-plan-interactive-ship-autonomous.md)

## History

- `01M67315J0N3C4HZEE7T2NNR55` Plan milestone 5: Rhiza usable from Rex, with an actionable roadmap app and an inbox app
