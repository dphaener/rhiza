# workflow-cleanup-plan-2026-11-11

source · `01M763Y5C2KHYR9MXYA5FVZ5SA`

## quotes

```json
[
  {
    "by": "owner",
    "text": "Placeholder only. The scaffold should have extensive comments for each item with multiple examples for popular languages so the developer can easily scan through and create the workflow without having to jump around in documentation.",
    "where": "2026-11-11 milestone plan interview"
  },
  {
    "by": "owner",
    "text": "We may want to chat about this, but isn't just reading the exit code good enough? If we need to detect what kind of error, maybe a quick  call out to an LLM (headless) using a model like Haiku would be a better choice. Let's talk about it.",
    "where": "2026-11-11 milestone plan interview"
  },
  {
    "by": "owner",
    "text": "Number 1, but we have to be aware that users can bring their own LLM so we need to leave that open. I think global configuration is a good place for things like this. ",
    "where": "2026-11-11 milestone plan interview"
  },
  {
    "by": "owner",
    "text": "I think that should all get cleaned up. Rhiza is the only self update candidate moving forward. Self update project is a useless setting at this point. Any non Rhiza projects that want self update can do so through their own custom workflow.",
    "where": "2026-11-11 milestone plan interview"
  },
  {
    "by": "owner",
    "text": "I think it can and should show up in the generated reference and user docs.",
    "where": "2026-11-11 milestone plan interview"
  },
  {
    "by": "owner",
    "text": "Self update SHOULD be a setting. I.E. if I am just a user of Rhiza, I should be able to turn on auto update and get the latest official release. It can also be just set to MAIN, which is what we would want here. For now I understand that main will always be the latest official release, but I will be adding a new release milestone soon that will change that so that the auto update can watch for actual releases and not just new builds.",
    "where": "2026-11-11 milestone plan interview"
  },
  {
    "by": "owner",
    "text": "How about we also add 'latest' or similar. We should have incremental builds available between builds that brave users can use. ",
    "where": "2026-11-11 milestone plan interview"
  },
  {
    "by": "owner",
    "text": "This should also be alled 'autoUpdate' and not 'selfUpdate', so it makes sense and aligns with the rest of the world",
    "where": "2026-11-11 milestone plan interview"
  },
  {
    "by": "owner",
    "text": "We can include all of the release channels in the plan, but it will rely on the release milestone. This planning needs to wrap up before we can begin the planning and execution of the release milestone.",
    "where": "2026-11-11 milestone plan interview"
  }
]
```

## raw_sha256

9f3155aed8beb5e2888e17b777b92bdffa6527f3ac0a59c0508a8569142a4f96

## session

70ee81c4-37df-450f-86ae-01a1cbc527f8

## summary

Milestone 32 planned with the owner. Measured on the VM: disk 79% used; the biggest use is Rhiza's own leftovers (Claude scratch folders under ~/tmp ~29-34G, index .db files for removed worktrees ~15G, stray rhiza test temp dirs), not the Go cache (go-build 1.4G, module cache 3.1G). Settled by picks: the product cleans up Rhiza's own leftovers at run end and the daemon prunes orphans; workflows gain an end-of-run cleanup phase that runs on done or failed and never fails the run; the Go cache trim (rhiza.daemon.trimGoCache, gosweep) and ship's head build move into Rhiza's own workflows; per-run GOTMPDIR becomes workflow environment and the step runner's Go-signature rerun uses the failure classifier; ci-triage hands the fixer the failed CI logs and the per-language parsers go; auto-update defaults to off; the 'release' channel and the 'latest' channel wait for the release milestone's published builds. Settled by typed answers: the workflow new scaffold is a placeholder with extensive comments and multi-language examples; failures are classified by a headless LLM (Haiku by default) configured globally so users can bring their own; self-update becomes a documented setting named autoUpdate with channels (off, main now; latest, release with the release milestone), self_update_project is removed. Seven change sets were opened and all planned now; the latest/release channels are blocked on the release milestone, which the owner plans after this one.

## title

Workflow cleanup and language-neutral defaults, planned with the owner

## Links

### In

- derived-from [d-scaffold-placeholder](../decision/d-scaffold-placeholder.md)
- derived-from [d-llm-classifier](../decision/d-llm-classifier.md)
- derived-from [d-no-self-update-project](../decision/d-no-self-update-project.md)
- derived-from [d-auto-update-documented](../decision/d-auto-update-documented.md)
- derived-from [d-auto-update-channels](../decision/d-auto-update-channels.md)
- derived-from [d-auto-update-latest](../decision/d-auto-update-latest.md)
- derived-from [d-auto-update-name](../decision/d-auto-update-name.md)
- derived-from [d-channels-planned-now](../decision/d-channels-planned-now.md)
- derived-from [answer-70ee81c4-1](../decision/answer-70ee81c4-1.md)
- derived-from [answer-70ee81c4-2](../decision/answer-70ee81c4-2.md)
- derived-from [answer-70ee81c4-3](../decision/answer-70ee81c4-3.md)
- derived-from [answer-70ee81c4-4](../decision/answer-70ee81c4-4.md)
- derived-from [answer-70ee81c4-5](../decision/answer-70ee81c4-5.md)
- derived-from [answer-70ee81c4-6](../decision/answer-70ee81c4-6.md)
- derived-from [answer-70ee81c4-7](../decision/answer-70ee81c4-7.md)
- derived-from [answer-70ee81c4-8](../decision/answer-70ee81c4-8.md)
- derived-from [answer-70ee81c4-9](../decision/answer-70ee81c4-9.md)
- derived-from [answer-70ee81c4-10](../decision/answer-70ee81c4-10.md)
- derived-from [answer-70ee81c4-11](../decision/answer-70ee81c4-11.md)
- derived-from [answer-70ee81c4-12](../decision/answer-70ee81c4-12.md)
- derived-from [answer-70ee81c4-13](../decision/answer-70ee81c4-13.md)
- derived-from [answer-70ee81c4-14](../decision/answer-70ee81c4-14.md)
- derived-from [answer-70ee81c4-15](../decision/answer-70ee81c4-15.md)
- derived-from [answer-70ee81c4-16](../decision/answer-70ee81c4-16.md)

## History

- `01M75W99SF6Q1T4QJM5GKA4PM1` Workflow cleanup and language-neutral defaults. Owner, 2026-11-11: "I'd like to propose a feature for rhiza. I think the Go cache is causing a lot of disk bloat, so the workflows should manage this when they are done. THere should be a cleanup phase. This is where the default Rhiza workflow (what is shipped and installed for all users as the default) from what we use to build Rhiza itself. This may actually be a milestone because I think it's going to cover a lot of ground. So we want to add a cleanup phase for THIS Rhiza only. But we may want to add some guidance in the docs, some workflow examples for different languages, etc. THis will require some milestone planning for sure. I also noticed when you were talking about the existing daemon features that the daemon was specifically looking for things like make check, which won't be universal to all programming languages, so things like that need to be workflow dependent, project dependent, or a global setting. Let's get to it!" Orchestrator's findings to start from (verify, don't trust): on the VM disk is 80% used (73G of 97G); ~/.cache/go-build 1.3G, ~/go/pkg/mod 3.1G, run worktrees 55-210M each. Go- and Rhiza-specific behaviour shipped in the default product today: the daemon setting rhiza.daemon.trimGoCache (internal/api/daemon.go:218, internal/gosweep) trims the Go build cache only when disk runs low; internal/tmpsweep sweeps go test temp folders; the watchdog's transient signatures include Go's go-build $WORK error (internal/api/daemon_watchdog.go:425); ship's head build runs go build ./cmd/rhiza only in Rhiza's own repo (internal/api/ship_build.go:85); daemon self-update runs make (internal/api/daemon.go:336); rhiza workflow new scaffolds run = "make test" (internal/api/workflow_new.go:75); ci-triage parses Go, RSpec, Minitest and pytest failures (internal/api/ship.go:134,1049); the local gate is already a project setting (rhiza.ship.localGate). The built-in workflows are internal/workflow/builtin/{plan,plan-milestone,build,review,ship-change}.toml; Rhiza's repo adds .rhiza/workflows/docs*.toml and the bugfix workflow (#435). Scope to settle with the owner: a cleanup step for Rhiza's own workflows (Go cache, worktrees, temp folders) versus the shipped defaults; which language- or tool-specific behaviours move into workflows, project settings or global settings; docs guidance and per-language workflow examples.
