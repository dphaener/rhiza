# m31-planning

source · `01M71G37GPD4RWWNN6JNBF8HJC`

## quotes

```json
[
  {
    "by": "owner",
    "text": "some of the examples are a little rough (using giant json blobs in cli commands, etc) and it doesn't look like it was all humanized using my voice",
    "where": "milestone request, 2026-11-10"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Plain commands: docs use rhiza add, rhiza link, rhiza set one-liners for every example; rhiza apply with a JSON file gets one short mention as the batch/agent path; no CLI change",
    "where": "interview, examples"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Owner sign-off per page: each doc change set's run still drafts with humanize, then stops for the owner's sign-off on the rendered prose before merge",
    "where": "interview, voice check"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Docs workflow: a first change set adds a repo workflow (ship-change with autonomy signoff = owner) and the doc change sets are run with it; other work stays unattended",
    "where": "interview, mechanism"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "README + quickstart, Concepts, Agents, runs, workflows, CLI help wording",
    "where": "interview, scope"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "All five now",
    "where": "interview, plan now"
  }
]
```

## raw_sha256

e83c97157606c453b39cb3640a805eb5f6e57c75ba9acbcc88c6845282e221b2

## session

97c939ee-4270-4def-bfd2-01a1c2e382d4

## summary

The owner asked for a documentation polish milestone: some doc examples are rough (giant JSON blobs in CLI commands) and the prose doesn't read in the owner's voice, although milestone 30's constraint doc-prose-humanized already required the humanize skill. Settled: doc examples use plain rhiza add/link/set commands, with rhiza apply mentioned once as the batch and agent path, and no CLI change for it. The voice is checked by the owner: each doc change set's run drafts with humanize and stops at an owner sign-off before merge. Mechanism: a repo workflow named docs (ship-change with signoff = owner), plus a way to pin a change set's workflow from rhiza change new and on an existing draft, since only rhiza request --workflow can today. Scope: README and quickstart, the concepts pages, the agents/runs/custom-workflows pages, and the CLI help wording behind the generated reference (plain and consistent, not in the owner's voice, normal unattended review). Five change sets, all planned now; the three page rewrites depend on the workflow change set. Nothing was left open.

## title

Milestone 31 planning: documentation polish, with the owner

## Links

### In

- derived-from [m31-polish-request](../decision/m31-polish-request.md)
- derived-from [doc-examples-plain-commands](../constraint/doc-examples-plain-commands.md)
- derived-from [answer-97c939ee-1](../decision/answer-97c939ee-1.md)
- derived-from [answer-97c939ee-2](../decision/answer-97c939ee-2.md)
- derived-from [answer-97c939ee-3](../decision/answer-97c939ee-3.md)
- derived-from [answer-97c939ee-4](../decision/answer-97c939ee-4.md)
- derived-from [answer-97c939ee-5](../decision/answer-97c939ee-5.md)

## History

- `01M71E6H2N1CYFN0F2NGAWWBG8` Documentation polish milestone for Rhiza. Owner, 2026-11-10: "some of the examples are a little rough (using giant json blobs in cli commands, etc) and it does not look like it was all humanized using my voice"
