# web-ui-milestone-interview

source · `01M6PFE13VHPKWY3QZ00HWMS8B`

## quotes

```json
[
  {
    "by": "owner",
    "text": "I want to build out the first pass of a web UI that can be used to monitor the progress of Rhiza. This will run as part of the daemon process and be served on localhost. The user will be responsible for any port forwarding (if it's running on this VM for instance). First pass will be read only. For monitoring purposes only and we will build from there.",
    "where": "the change set's intent, as filed"
  },
  {
    "by": "owner",
    "text": "The first requirement will be to obtain a design. You will generate a prompt that I will then feed to Claude design to generate a design system.",
    "where": "the change set's intent, as filed"
  },
  {
    "by": "owner",
    "text": "I already have some design elements to work with.",
    "where": "plan session, before the interview"
  },
  {
    "by": "owner",
    "text": "Rhiza — “Bioluminescent soil.” Dark-first, calm, precise.",
    "where": "plan session, the design basics the owner pasted"
  },
  {
    "by": "owner",
    "text": "Motifs: fine glowing hyphae threads, small typed nodes (filled dot = intent, ring = capability, diamond = decision, square = constraint), always-dark terminals, pill buttons, 12–16px radii, thin 1px rules. No gradients, no emoji.",
    "where": "plan session, the design basics the owner pasted"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Go templates, no build; internal/api read ops; Server-sent events; On by default, 127.0.0.1; Overview, Runs, Roadmap & change sets; Prompt now, design change set waits; All five now; 7787; Fonts embedded in the binary",
    "where": "plan session, AskUserQuestion picks"
  }
]
```

## raw_sha256

c6bc38b2ddb590dd3012274b9980b0599485bea05228ec8898ce620bd972dd88

## session

ea442c0d-1f6e-46a3-908f-01a1acf1a3a2

## summary

Plan session (2026-11-05) for change set 01M6PF31, a record for a new milestone: a first-pass web UI for monitoring Rhiza, run inside the daemon on localhost, read-only, with the owner forwarding the port. The owner gave the design basics (palette, type, motifs) and asked for a Claude Design prompt; the agent wrote it in the session for the owner to feed to Claude Design. Settled through the agent's options, each picked by the owner: Go html/template pages with no build step, embedded with go:embed; data from internal/api read operations through a read-only allowlist, also exposed as JSON at /api/<op>; live updates by server-sent events; served on by default on 127.0.0.1, port 7787, changeable or off via rhiza daemon set rhiza.web.port; fonts embedded in the binary; first-pass pages are Overview, Runs, and Roadmap & change sets (the model browser is left for a later pass); the design system change set waits on a question until Claude Design's output arrives, and the pages depend on it while the server does not; all five change sets (server, design system, Overview, Runs, Roadmap & change sets) planned now. Open: Claude Design's output, which blocks the design system leaves.

## title

Milestone 28 planning: the first read-only web UI, with the owner

## Links

### In

- derived-from [web-ui-in-daemon-read-only](../decision/web-ui-in-daemon-read-only.md)
- derived-from [web-ui-design-from-claude-design](../decision/web-ui-design-from-claude-design.md)
- derived-from [web-ui-design-basics](../decision/web-ui-design-basics.md)
- derived-from [answer-ea442c0d-1](../decision/answer-ea442c0d-1.md)
- derived-from [answer-ea442c0d-2](../decision/answer-ea442c0d-2.md)
- derived-from [answer-ea442c0d-3](../decision/answer-ea442c0d-3.md)
- derived-from [answer-ea442c0d-4](../decision/answer-ea442c0d-4.md)
- derived-from [answer-ea442c0d-5](../decision/answer-ea442c0d-5.md)
- derived-from [answer-ea442c0d-6](../decision/answer-ea442c0d-6.md)
- derived-from [answer-ea442c0d-7](../decision/answer-ea442c0d-7.md)
- derived-from [answer-ea442c0d-8](../decision/answer-ea442c0d-8.md)
- derived-from [answer-ea442c0d-9](../decision/answer-ea442c0d-9.md)

## History

- `01M6PF31YCJ041HB94MDWQR331` I want to build out the first pass of a web UI that can be used to monitor the progress of Rhiza. This will run as part of the daemon process and be served on localhost. The user will be responsible for any port forwarding (if it's running on this VM for instance). First pass will be read only. For monitoring purposes only and we will build from there. The first requirement will be to obtain a design. You will generate a prompt that I will then feed to Claude design to generate a design system.
