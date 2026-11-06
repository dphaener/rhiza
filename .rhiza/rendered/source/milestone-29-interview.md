# milestone-29-interview

source · `01M6RN1GT6AJ10J39S4RTN28YS`

## quotes

```json
[
  {
    "by": "owner",
    "text": "There isn't anything existing so root can be required from the start",
    "where": "answer on what happens to a repo under no root"
  },
  {
    "by": "owner",
    "text": "The TUI needs to become project aware, so when I'm starting a new milestone/changeset or whatever, it needs to know which project I am currently in the context of.",
    "where": "message during the interview"
  }
]
```

## raw_sha256

eadb1cadac63c9ac0fa9a56ee991ee87837a8f9fc71d7ac68e1d85f7d5595e2d

## session

05fe56c8-8c19-4241-83c6-01a1b1494bc7

## summary

The owner planned milestone 29, in which Rhiza manages and builds several projects on one machine from one root directory. Picked options: one root daemon per machine drives every project under one usage governor and one machine-wide run cap, and serves one web UI. Projects are registered explicitly under the root (rhiza project add, by path or git URL). Run slots are shared round-robin across projects, in each project's roadmap order, with an optional per-project max. Machine-wide settings live in the root config and per-project settings stay in each repo's git config. The web UI gets an all-projects home and a merged inbox, plus per-project pages through a switcher. The CLI and the Rex TUI apps go multi-project in this milestone. Worktrees stay as <root>/<project>.<slug> siblings. A registered project with no .rhiza/ is listed but not built. Each project keeps its own Rex session. All six change sets (root and registry, root daemon, cross-project CLI, web UI, TUI apps, VM migration) are planned now. In a typed answer, the owner ruled that a root is required from the start, with no per-repo daemon mode kept. Mid-interview the owner added that the TUI must be project aware: starting a milestone, change set or anything else happens in the context of the project the owner is currently in.

## title

Milestone 29 interview: many projects from one root

## Links

### In

- derived-from [m29-one-root-daemon](../decision/m29-one-root-daemon.md)
- derived-from [m29-registered-projects](../decision/m29-registered-projects.md)
- derived-from [m29-round-robin-slots](../decision/m29-round-robin-slots.md)
- derived-from [m29-settings-split](../decision/m29-settings-split.md)
- derived-from [m29-web-all-projects](../decision/m29-web-all-projects.md)
- derived-from [m29-cli-and-tui-now](../decision/m29-cli-and-tui-now.md)
- derived-from [m29-plan-all-now](../decision/m29-plan-all-now.md)
- derived-from [m29-worktrees-siblings](../decision/m29-worktrees-siblings.md)
- derived-from [m29-root-required](../decision/m29-root-required.md)
- derived-from [m29-unready-listed](../decision/m29-unready-listed.md)
- derived-from [m29-rex-session-per-project](../decision/m29-rex-session-per-project.md)
- derived-from [m29-tui-current-project](../decision/m29-tui-current-project.md)
- derived-from [answer-05fe56c8-12](../decision/answer-05fe56c8-12.md)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
