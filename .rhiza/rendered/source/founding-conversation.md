# founding-conversation

source · `01M5SWRGJ7AZM5CYJAE32SGRYM`

## quotes

```json
[
  {
    "by": "owner",
    "text": "Ok we're going to start planning something ambitious. We're running in the just released Rex\nmultiplexer and it has a lot of great features. I already have something I've built called afk\n(~/code/afk). I want to build on those principles but start something new, completely from scratch,\nfrom the ground up. Here is the high level:\n\nA platform, much like afk, that is built for all levels of developers, product managers, ux\ndesigners, hobbyists, everyone. You can use as little of it as you want, bring your own workflows,\nwhatever coding agent you want, it also has all the bells and whistles (shipit/forge workflow,\nopinionated plan/build workflows) built in. Hobbyists can tinker with it. Casual users can just use\nthe web ui and still build something useful.\n\nI only have access to the MacOS trial right now, and I'm a hard core TUI user, and super developer\nnerd so we're going to start with core building blocks and deep dev features. Think core APIs, CLIs,\netc. Things that are not friendly to most users, but are the building blocks to the later features.\n\nSimultaneously, I want to be building out the new engine that drives feature and product\ndevelopment. The docs/PLAN.md was a good start but it grows unwieldy quickly, doesn't really\nrepresent what a product is. A living, breathing, ever changing piece of software, that has records\nof how it has changed, how the decisions were made, and how all of it is connected. More of a\ninterconnected web than a linear page of data. How we describe the product, store it, and morph it\nis the core, and the system is built around that.\n\nDo not rush to start writing docs and making decisions, we have a LOT to discuss here. I want to\nbuild a sophisticated system that I can start using soon, that is going to be quality software, and\nthat will build quality software.",
    "where": "Opening message: the founding vision (## Owner, line 11)"
  },
  {
    "by": "owner",
    "text": "And let's not forget to record this initial conversation, it's going to be important. 😄",
    "where": "Opening, sent while the assistant read afk: record this conversation (## Owner, line 39)"
  },
  {
    "by": "owner",
    "text": "read this too: https://www.anthropic.com/research/redesigning-source-control - Then let's step\nthrough the questions one by one.",
    "where": "Before A1: step through the questions one by one (## Owner, line 117)"
  },
  {
    "by": "owner",
    "text": "1. I'm not sure. Initially we need to serve engineers, but also PMs at some point...\n2. Agreed, tasks belong to the work layer and are disposable\n3. Agreed, but custom kinds are down the roadmap\n4. I think pretty atomic, to keep context retrieval small and accurate, unless that's not a good\nidea? I think nodes could potentially have several sizes. I.E. a feature description might have\nmany smaller nodes that represent it's capabilities?",
    "where": "A1 reply: the spine, tasks, custom kinds, node size (## Owner, line 211)"
  },
  {
    "by": "owner",
    "text": "yes to both, let's move on to A2",
    "where": "A1 confirmation: one tree of capabilities, a feature is a capability with children (## Owner, line 272)"
  },
  {
    "by": "owner",
    "text": "One quick interjection that I forgot to include, I think part of the initial \"power user\"\nexperience should follow the Rex model and include a CLI and ship with skills that tell agents how\nto use it.",
    "where": "A2 interjection: a CLI that ships skills telling agents how to use it (## Owner, line 349)"
  },
  {
    "by": "owner",
    "text": "yes to all four, move on to A3",
    "where": "A2 reply: overlays, land at ship, land whole, intention and record (## Owner, line 385)"
  },
  {
    "by": "owner",
    "text": "1. This has to be portable to many repositories and live inside of them. I.E. if I'm using it at\nwork, I have to work across 15-20 different repos that are all different products with different\nroadmaps. I don't think a database is the right choice here. I think an append-only log is the\nright choice. The SQLite index is an interesting idea, is that agent friendly and performant?\n2. Yes it must, as I just described. 😄\n3. Yes\n4. Yes, but I think raw conversations like this can be summarized and scrubbed by an agent, with a\nshipped skill.\n\nYes I think rendered markdown should be comitted, it's useful and could prove useful in a future\nUI, a CI check is great, installable from the CLI or by an agent with a skill/CLI. Sources live in\nthe repo right? Didn't we already answer that? Or is this a different source?",
    "where": "A3 reply: append-only log, in each repo, IDs and handles, CLI-only writes, committed markdown (## Owner, line 467)"
  },
  {
    "by": "owner",
    "text": "I think it's easy enough to just always use a .<whatever_the_name_is> folder for storage inside the\nrepo and just global gitignore it if you don't want it comitted.",
    "where": "A3 follow-up: personal mode (always a dot-folder in the repo) (## Owner, line 553)"
  },
  {
    "by": "owner",
    "text": "yes to all four, move on to A5",
    "where": "A4 reply: staying true to the code (## Owner, line 622)"
  },
  {
    "by": "owner",
    "text": "yes to all four, move on to B6",
    "where": "A5 reply: repos and products (## Owner, line 686)"
  },
  {
    "by": "owner",
    "text": "yes to all four, move on to B7",
    "where": "B6 reply: the six layers (## Owner, line 749)"
  },
  {
    "by": "owner",
    "text": "1. Rex is required\n2. tmux is not an option\n3. Yes\n4. Yes",
    "where": "B7 reply: Rex required, no tmux (## Owner, line 815)"
  },
  {
    "by": "owner",
    "text": "yes to all four, move on to C",
    "where": "B8 reply: authorship and permissions (## Owner, line 876)"
  },
  {
    "by": "owner",
    "text": "1. Yes\n2. M17 can continue\n3. Yes\n4. Let's workshop this - you saved this conversation in a folder called genesis, which was a nice\nname I thought.\n5. Yes!",
    "where": "C9-C11 reply: clean repo, afk M17, Go, the name, how milestone 1 is built (## Owner, line 959)"
  },
  {
    "by": "owner",
    "text": "It was origin, but lineage makes more sense. Genesis just sounds cool. lol",
    "where": "Naming: why the name genesis (## Owner, line 1020)"
  },
  {
    "by": "owner",
    "text": "yes, Genesis it is, move on to the log format",
    "where": "Naming: Genesis it is (## Owner, line 1042)"
  },
  {
    "by": "owner",
    "text": "yes to all five, move on to the CLI commands",
    "where": "The log on disk (## Owner, line 1151)"
  },
  {
    "by": "owner",
    "text": "yes to all five, move on to the slice format",
    "where": "The CLI's commands (## Owner, line 1221)"
  },
  {
    "by": "owner",
    "text": "yes to all five, move on to the agent guide",
    "where": "The slice format (## Owner, line 1335)"
  },
  {
    "by": "owner",
    "text": "yes to all five, move on to the next topic",
    "where": "The agent guide (## Owner, line 1418)"
  },
  {
    "by": "owner",
    "text": "yes to all five, move on to the workflow engine",
    "where": "How a change set becomes work (## Owner, line 1500)"
  },
  {
    "by": "owner",
    "text": "yes to all five, move on to the agent adapter. Genesis should also be able to tell an agent how to\nwrite a custom workflow.",
    "where": "The workflow engine, and writing custom workflows (## Owner, line 1625)"
  },
  {
    "by": "owner",
    "text": "yes to all six, move on to where agents run",
    "where": "The agent adapter (## Owner, line 1722)"
  },
  {
    "by": "owner",
    "text": "yes to all four, let's cover the remaining topics first",
    "where": "Where agents run (## Owner, line 1794)"
  },
  {
    "by": "owner",
    "text": "yes to all four, move on to testing",
    "where": "Remaining topic 1: the link types (## Owner, line 1864)"
  },
  {
    "by": "owner",
    "text": "yes to all six, move on to the capability tree",
    "where": "Remaining topic 2: testing Genesis itself (## Owner, line 1944)"
  },
  {
    "by": "owner",
    "text": "yes to all four, move on to the web UI",
    "where": "Remaining topic 3: the first capability tree (## Owner, line 2033)"
  },
  {
    "by": "owner",
    "text": "yes to all seven, move on to the review workflow",
    "where": "Remaining topic 4: what the web UI will need (## Owner, line 2103)"
  },
  {
    "by": "owner",
    "text": "yes to all seven, move on to milestone 1 scope",
    "where": "Remaining topic 5: the review workflow (## Owner, line 2189)"
  },
  {
    "by": "owner",
    "text": "yes to all five, create the repo",
    "where": "Milestone 1 scope (## Owner, line 2276)"
  },
  {
    "by": "owner",
    "text": "Please listen for the tasks being done, you're the orchestrator here.",
    "where": "orchestration: 'you're the orchestrator here' (## Owner, line 2511)"
  },
  {
    "by": "owner",
    "text": "Make sure to ensure that eveyrthing is fully close. Branches/worktrees gone, CI green, and new binary shipped.",
    "where": "orchestration: when a task is fully closed (## Owner, line 2527)"
  },
  {
    "by": "owner",
    "text": "Can't you observe the event stream from the Rex.api for those tabs and just react?",
    "where": "orchestration: react to the Rex event stream (## Owner, line 2542)"
  },
  {
    "by": "owner",
    "text": "You're supposed to perform the merge once CI is green",
    "where": "orchestration: 'You're supposed to perform the merge' (## Owner, line 2570)"
  },
  {
    "by": "owner",
    "text": "This is great. I love the reviews. If you feel you need to interject with the decisions I made and update a plan, please do, and communicate with the other agents.",
    "where": "orchestration: the orchestrator may amend plans (## Owner, line 2595)"
  },
  {
    "by": "owner",
    "text": "This needs to be a feature in Genesis.",
    "where": "orchestration: 'This needs to be a feature in Genesis' (## Owner, line 2620)"
  },
  {
    "by": "owner",
    "text": "Milestone 2 sounds fine.",
    "where": "orchestration: leaf 6 is milestone 2's first change set (## Owner, line 2651)"
  },
  {
    "by": "owner",
    "text": "Those are from the ask user question tool, so none of them are my own words.",
    "where": "orchestration: picked options are not the owner's words (## Owner, line 2662)"
  },
  {
    "by": "owner",
    "text": "Let's do it",
    "where": "Task 11: go-ahead to start the import (## Owner, line 2802)"
  }
]
```

## raw_sha256

03520c0939d3b0f415ecbdce88c838e3eb541a731edc75c533a22e3ba7f5dd6c

## summary

The founding conversation of Genesis (from 2026-10-23), between the owner and Claude Code in Rex, later continued as the orchestrator's log of milestone 1. The owner set out to build, from scratch on the principles of an earlier tool of theirs (afk), a platform on which developers, product managers, designers and hobbyists build quality software with agents, using as little or as much of it as they want, with any agent and any workflow. It starts with building blocks for power users (APIs, CLIs); at its centre is an engine that describes a product as a living, linked record of what it is, how it changed and why, in place of a linear plan document. Each topic was settled by an owner message; the "Settled for now" summaries are the assistant's wording.

Decided. Model (A1): a small fixed kernel of node kinds with typed links; custom kinds later on the roadmap; tasks belong to a disposable work layer; one containment tree of capabilities, with components and journeys as views; a feature is a capability with children; nodes are atomic and reading is composed into slices. Change (A2): overlays on one as-is, with to-be and what-if computed; change sets land at ship and land whole, larger work being a sequence; two kinds, intention and record. A CLI with an agent guide shipped inside the binary is part of the first power-user experience. Storage (A3): an append-only operation log in text files inside each product's repo, working across many repos; random IDs plus renamable handles; only the CLI writes the log, enforced by check in CI; rendered markdown is committed and kept current by a CI check; a derived SQLite index; Source nodes hold a scrubbed summary and a hash, and raw transcripts never enter the repo; always a dot-folder in the repo, gitignored if wanted, with no separate personal mode. Truth to code (A4): checkable anchors, Change: commit trailers, reconcile only on deterministic triggers. Repos (A5): one model per repo with several roots for a monorepo; no home repo; cross-repo links by repo ID through a local registry; shared constraints from standards repos and a personal layer. Kernel (B6-B8): six layers (model store, model tools, runtime, work, orchestration, surfaces), each usable alone; the CLI works on files with no daemon; layers 1-2 first, first used to import this conversation; hooks on model events. Rex is required and there is no tmux fallback (the owner chose this against the assistant's recommendation); state never lives in Rex. Every entry records an actor and who it acted for; actors are nodes; no permission enforcement at first, but accepts and rejects record their rule; no hidden nodes. Practicalities (C9-C11): a clean repo with afk pieces ported on purpose; afk's M17 continues; Go, one binary that is also the daemon, pure domain packages under a thin CLI; milestone 1 built with Claude Code in Rex and shipit. The name is Genesis (for lineage, and because it sounds cool): command genesis, folder .genesis/, no official short command.

Design topics then settled: the log on disk (one immutable file per entry, the envelope, six operations with field-level updates against base versions, causal order, canonical versioned JSON, never compacted); the CLI's shape and rules (a current change set per worktree, --record, apply, why, --json, --dry-run, stable exit codes, next in errors, no prompts without a terminal); the slice (ten sections, inheritance down the tree, a 4,000-token budget with fixed cuts, the same bytes every time); the agent guide (three layers, three skills, a size cap); the work boundary (one change set per PR, roadmap and ready as views, the write-back contract); the workflow engine (declarative TOML state machines run by the engine, and, added by the owner, teaching agents to write custom workflows); the agent adapter (Claude Code first, GENESIS_ACTOR); where agents run (supervised on the Mac until layer 5, bypassed permissions only on dedicated machines); the fourteen link types; testing; the first capability tree (twelve areas, four root intents, seven root constraints, the shape of the import); web-UI readiness; the review workflow. Milestone 1: scope in and out, five done points, eleven tasks, a temporary task list to be retired by the import, and a private GitHub repo.

The later part records milestone 1's orchestration. The owner made this session the orchestrator: it watches the task sessions, merges on green CI, closes each task out fully (branch and worktree gone, CI green on main, a new binary installed), and may amend a task's plan when it conflicts with the founding decisions or another task, telling the affected sessions. The owner asked that this orchestrated supervision become a Genesis feature: six proposed leaves under the orchestration area, with leaf 6 (interview answers as Decision nodes) as milestone 2's first change set. The owner settled that an option picked from an agent's question menu is not the owner's words. Orchestrator notes record decisions it made under the owner's authorization (one default-view rule for reads and slice, a 50 ms read budget in CI) and follow-ups for the import: incremental what-if, a check warning for unlinked constraints, a field schema per kind, cursor paging on log, roadmap and why, faster roadmap, ready and index reads, the orchestration leaves, writes that see accepted unlanded dependencies, and link checking endpoint kinds.

Open: acceptance policy (layer 5); journeys and other PM-facing views; custom kinds and link types; where unattended agents run at layer 5 (a Linux hub if Rex reaches it, otherwise a macOS VM spike); everything listed as out of milestone 1.

## title

The founding conversation, with the owner

## Links

### In

- derived-from [genesis](../capability/genesis.md)
- derived-from [model](../capability/genesis/model.md)
- derived-from [change-sets](../capability/genesis/change-sets.md)
- derived-from [lineage](../capability/genesis/lineage.md)
- derived-from [views](../capability/genesis/views.md)
- derived-from [integrity](../capability/genesis/integrity.md)
- derived-from [guide](../capability/genesis/guide.md)
- derived-from [anyone-builds](../intent/anyone-builds.md)
- derived-from [model-tells-truth](../intent/model-tells-truth.md)
- derived-from [small-accurate-context](../intent/small-accurate-context.md)
- derived-from [power-users-first](../intent/power-users-first.md)
- derived-from [log-append-only-cli-writes](../constraint/log-append-only-cli-writes.md)
- derived-from [same-model-same-output](../constraint/same-model-same-output.md)
- derived-from [agents-named-only-in-adapters](../constraint/agents-named-only-in-adapters.md)
- derived-from [state-never-in-rex](../constraint/state-never-in-rex.md)
- derived-from [rex-required-from-runtime](../constraint/rex-required-from-runtime.md)
- derived-from [no-prompts-without-tty](../constraint/no-prompts-without-tty.md)
- derived-from [bypass-only-on-dedicated](../constraint/bypass-only-on-dedicated.md)
- derived-from [tasks-are-disposable](../decision/tasks-are-disposable.md)
- derived-from [fixed-core-kinds](../decision/fixed-core-kinds.md)
- derived-from [one-capability-tree](../decision/one-capability-tree.md)
- derived-from [feature-is-capability](../decision/feature-is-capability.md)
- derived-from [overlays-on-as-is](../decision/overlays-on-as-is.md)
- derived-from [land-at-ship](../decision/land-at-ship.md)
- derived-from [land-whole](../decision/land-whole.md)
- derived-from [two-change-set-kinds](../decision/two-change-set-kinds.md)
- derived-from [cli-with-agent-guide](../decision/cli-with-agent-guide.md)
- derived-from [append-only-log](../decision/append-only-log.md)
- derived-from [log-in-product-repo](../decision/log-in-product-repo.md)
- derived-from [random-ids-and-handles](../decision/random-ids-and-handles.md)
- derived-from [only-cli-writes-log](../decision/only-cli-writes-log.md)
- derived-from [rendered-markdown-committed](../decision/rendered-markdown-committed.md)
- derived-from [always-a-dot-folder](../decision/always-a-dot-folder.md)
- derived-from [anchor-shape](../decision/anchor-shape.md)
- derived-from [change-trailers](../decision/change-trailers.md)
- derived-from [reconcile-on-trigger](../decision/reconcile-on-trigger.md)
- derived-from [ci-fails-on](../decision/ci-fails-on.md)
- derived-from [ci-warns-on](../decision/ci-warns-on.md)
- derived-from [one-model-per-repo](../decision/one-model-per-repo.md)
- derived-from [no-home-repo](../decision/no-home-repo.md)
- derived-from [links-by-repo-id](../decision/links-by-repo-id.md)
- derived-from [standards-and-personal-layer](../decision/standards-and-personal-layer.md)
- derived-from [six-layers](../decision/six-layers.md)
- derived-from [cli-needs-no-daemon](../decision/cli-needs-no-daemon.md)
- derived-from [layers-1-2-first](../decision/layers-1-2-first.md)
- derived-from [hooks-from-the-start](../decision/hooks-from-the-start.md)
- derived-from [rex-required](../decision/rex-required.md)
- derived-from [no-tmux-fallback](../decision/no-tmux-fallback.md)
- derived-from [rex-holds-no-state](../decision/rex-holds-no-state.md)
- derived-from [rex-extras-optional](../decision/rex-extras-optional.md)
- derived-from [actor-on-every-entry](../decision/actor-on-every-entry.md)
- derived-from [actors-are-nodes](../decision/actors-are-nodes.md)
- derived-from [no-enforcement-at-first](../decision/no-enforcement-at-first.md)
- derived-from [no-hidden-nodes](../decision/no-hidden-nodes.md)
- derived-from [clean-repo](../decision/clean-repo.md)
- derived-from [afk-m17-continues](../decision/afk-m17-continues.md)
- derived-from [go-one-binary](../decision/go-one-binary.md)
- derived-from [built-with-claude-code](../decision/built-with-claude-code.md)
- derived-from [name-is-genesis](../decision/name-is-genesis.md)
- derived-from [file-per-entry](../decision/file-per-entry.md)
- derived-from [entry-envelope](../decision/entry-envelope.md)
- derived-from [six-operation-kinds](../decision/six-operation-kinds.md)
- derived-from [causal-order](../decision/causal-order.md)
- derived-from [canonical-json](../decision/canonical-json.md)
- derived-from [rex-command-shape](../decision/rex-command-shape.md)
- derived-from [change-set-per-worktree](../decision/change-set-per-worktree.md)
- derived-from [apply-is-main-write](../decision/apply-is-main-write.md)
- derived-from [why-is-lineage-verb](../decision/why-is-lineage-verb.md)
- derived-from [cli-rules](../decision/cli-rules.md)
- derived-from [ten-slice-sections](../decision/ten-slice-sections.md)
- derived-from [inherited-down-tree](../decision/inherited-down-tree.md)
- derived-from [slice-token-budget](../decision/slice-token-budget.md)
- derived-from [same-model-same-bytes](../decision/same-model-same-bytes.md)
- derived-from [slice-change-way-in](../decision/slice-change-way-in.md)
- derived-from [three-guide-layers](../decision/three-guide-layers.md)
- derived-from [skill-contents](../decision/skill-contents.md)
- derived-from [three-skills](../decision/three-skills.md)
- derived-from [user-level-plugin](../decision/user-level-plugin.md)
- derived-from [size-cap-and-evals](../decision/size-cap-and-evals.md)
- derived-from [one-change-set-one-pr](../decision/one-change-set-one-pr.md)
- derived-from [no-milestone-kind](../decision/no-milestone-kind.md)
- derived-from [ready-pure-function](../decision/ready-pure-function.md)
- derived-from [work-contract](../decision/work-contract.md)
- derived-from [work-state-outside-log](../decision/work-state-outside-log.md)
- derived-from [engine-orchestrates](../decision/engine-orchestrates.md)
- derived-from [five-step-kinds](../decision/five-step-kinds.md)
- derived-from [declarative-toml](../decision/declarative-toml.md)
- derived-from [run-journal](../decision/run-journal.md)
- derived-from [built-in-workflows-by-name](../decision/built-in-workflows-by-name.md)
- derived-from [explain-custom-workflows](../decision/explain-custom-workflows.md)
- derived-from [workflow-writing-tools](../decision/workflow-writing-tools.md)
- derived-from [adapters-declare-capabilities](../decision/adapters-declare-capabilities.md)
- derived-from [built-in-go-adapters](../decision/built-in-go-adapters.md)
- derived-from [status-from-hooks-first](../decision/status-from-hooks-first.md)
- derived-from [claude-code-adapter](../decision/claude-code-adapter.md)
- derived-from [actor-set-at-launch](../decision/actor-set-at-launch.md)
- derived-from [supervised-on-the-mac](../decision/supervised-on-the-mac.md)
- derived-from [trust-per-machine](../decision/trust-per-machine.md)
- derived-from [rex-host-by-name](../decision/rex-host-by-name.md)
- derived-from [linux-hub-at-layer-5](../decision/linux-hub-at-layer-5.md)
- derived-from [parent-and-authorship](../decision/parent-and-authorship.md)
- derived-from [fourteen-links](../decision/fourteen-links.md)
- derived-from [no-proposes-link](../decision/no-proposes-link.md)
- derived-from [fixed-link-types](../decision/fixed-link-types.md)
- derived-from [table-and-property-tests](../decision/table-and-property-tests.md)
- derived-from [golden-files](../decision/golden-files.md)
- derived-from [testscript-scenarios](../decision/testscript-scenarios.md)
- derived-from [performance-budget](../decision/performance-budget.md)
- derived-from [permanent-corpus](../decision/permanent-corpus.md)
- derived-from [checks-itself-in-ci](../decision/checks-itself-in-ci.md)
- derived-from [tree-conventions](../decision/tree-conventions.md)
- derived-from [twelve-areas](../decision/twelve-areas.md)
- derived-from [root-intents-and-constraints](../decision/root-intents-and-constraints.md)
- derived-from [import-shape](../decision/import-shape.md)
- derived-from [transport-independent-ops](../decision/transport-independent-ops.md)
- derived-from [change-feed-events](../decision/change-feed-events.md)
- derived-from [references-by-id](../decision/references-by-id.md)
- derived-from [field-level-diffs](../decision/field-level-diffs.md)
- derived-from [notes-are-discussion](../decision/notes-are-discussion.md)
- derived-from [cursor-paged-lists](../decision/cursor-paged-lists.md)
- derived-from [explicit-actor-and-change](../decision/explicit-actor-and-change.md)
- derived-from [code-checks-first](../decision/code-checks-first.md)
- derived-from [three-agent-checks](../decision/three-agent-checks.md)
- derived-from [critics-on-serious-findings](../decision/critics-on-serious-findings.md)
- derived-from [re-review-failures-only](../decision/re-review-failures-only.md)
- derived-from [review-profiles-by-rule](../decision/review-profiles-by-rule.md)
- derived-from [review-feeds-model](../decision/review-feeds-model.md)
- derived-from [one-pr-summary](../decision/one-pr-summary.md)
- derived-from [milestone-1-scope](../decision/milestone-1-scope.md)
- derived-from [five-done-points](../decision/five-done-points.md)
- derived-from [eleven-tasks](../decision/eleven-tasks.md)
- derived-from [temporary-task-list](../decision/temporary-task-list.md)
- derived-from [repo-location](../decision/repo-location.md)
- derived-from [private-github-repo](../decision/private-github-repo.md)
- derived-from [session-orchestrates](../decision/session-orchestrates.md)
- derived-from [task-close-out](../decision/task-close-out.md)
- derived-from [orchestrator-merges](../decision/orchestrator-merges.md)
- derived-from [orchestrator-amends-plans](../decision/orchestrator-amends-plans.md)
- derived-from [supervision-is-a-capability](../decision/supervision-is-a-capability.md)
- derived-from [leaf-6-in-milestone-2](../decision/leaf-6-in-milestone-2.md)
- derived-from [picked-options-not-owner-words](../decision/picked-options-not-owner-words.md)
- derived-from [t5-d13-budget-threshold](../decision/t5-d13-budget-threshold.md)
- derived-from [t7-d2-default-view](../decision/t7-d2-default-view.md)
- derived-from [t7-d5-why-walk](../decision/t7-d5-why-walk.md)
- derived-from [t7-d14-read-budget](../decision/t7-d14-read-budget.md)
- derived-from [t8-d1-slice-shape](../decision/t8-d1-slice-shape.md)
- derived-from [t10-fix1-only-typed-words](../decision/t10-fix1-only-typed-words.md)
- derived-from [t10-fix2-work-contract](../decision/t10-fix2-work-contract.md)
- derived-from [t10-fix3-minor-review-points](../decision/t10-fix3-minor-review-points.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
