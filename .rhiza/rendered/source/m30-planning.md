# m30-planning

source · `01M6W5ANB6CXAAYAVSP53Z43CY`

## quotes

```json
[
  {
    "by": "owner",
    "text": "i want to start a new milestone first. I'm planning on sharing this out with few people soon and i need to get a polished readme and documentation in place. if my humanize plugin isn't installed we need to install it for all the docs. it's in my dotfiles",
    "where": "planning chat, 2026-11-07"
  }
]
```

## raw_sha256

a34bb6c0bbacc3adca0cb8cf388cd5955381bdb6c896103669217955bd86010a

## summary

The owner plans to share Rhiza with a few people soon and wants a polished README and documentation. The humanize skill (voice engine from the owner's dotfiles) was installed first so all doc prose goes through it. Settled through multiple-choice picks of the agent's options: the audience is developers who try Rhiza on their own repo; install by go install plus tagged GitHub releases built with GoReleaser, first version v0.1.0; docs are Markdown in docs/ with the README as front door (the tabled landing page stays tabled); docs lead with the CLI model working alone, with Rex and Claude Code as requirements for runs and the daemon; the doc set is README and quickstart, concepts, agents/runs/daemon, and reference; the CLI, config and workflow reference is generated and checked stale by make check; the repo goes public under MIT, after an audit of git history and .rhiza/ for secrets, personal data and private hostnames, with the owner flipping visibility by hand; humanize is a standing rule for human-facing doc prose, generated reference exempt; the quickstart covers the model loop without agents so a testscript can run it; all six change sets are planned now; the agents/runs/daemon docs wait for the queued workflow work (dispatch, per-workflow autonomy, schedules, refs, warnings) to land.

## title

Milestone 30 planning: ready to share

## Links

### In

- derived-from [doc-prose-humanized](../constraint/doc-prose-humanized.md)
- derived-from [m30-audience](../decision/m30-audience.md)
- derived-from [m30-install-path](../decision/m30-install-path.md)
- derived-from [m30-docs-in-repo](../decision/m30-docs-in-repo.md)
- derived-from [m30-core-alone](../decision/m30-core-alone.md)
- derived-from [m30-doc-set](../decision/m30-doc-set.md)
- derived-from [m30-generated-reference](../decision/m30-generated-reference.md)
- derived-from [m30-go-public](../decision/m30-go-public.md)
- derived-from [m30-mit](../decision/m30-mit.md)
- derived-from [m30-audit-then-owner-flips](../decision/m30-audit-then-owner-flips.md)
- derived-from [m30-humanize-rule](../decision/m30-humanize-rule.md)
- derived-from [m30-goreleaser](../decision/m30-goreleaser.md)
- derived-from [m30-v0-1-0](../decision/m30-v0-1-0.md)
- derived-from [m30-quickstart-model-loop](../decision/m30-quickstart-model-loop.md)
- derived-from [m30-plan-all-now](../decision/m30-plan-all-now.md)
- derived-from [m30-agent-docs-after-workflows](../decision/m30-agent-docs-after-workflows.md)

## History

- `01M6W5A2R4BXDGQGHNWGZNB3XY` Milestone 30: Rhiza is ready to share: install, README and docs, generated reference, public under MIT
