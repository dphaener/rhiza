# plan-web-roadmap-project-scope

source · `01M79S98X7E9Y23AJREFNTD180`

## quotes

```json
[
  {
    "by": "owner",
    "text": "The roadmap tab doesn't filter between projects properly, it shows things from both projects",
    "where": "bug request, 2026-11-13"
  }
]
```

## raw_sha256

8ede7393996eda46094dbd18ac4c5611a18a1f8f38a352c76f37b0cc4b3a3957

## session

e85a7532-fc4a-4bf1-806f-01a1d392d448

## summary

Autonomous plan run (no owner interview) for the bug 'The roadmap tab doesn't filter between projects properly'. Reading the code: internal/web/roadmap.go loadRoadmap and internal/web/change.go dependsOn read /api/roadmap?all=true through the project's ReadOnlyHandler (internal/api/web_api.go). Under a Rhiza root, runRoadmap (internal/api/plan.go) passes all=true to inboxScopes, which widens to every registered project, so /p/<project>/roadmap lists every project's change sets and milestones. health and web_milestones avoid this with the unexported local flag, which an HTTP query cannot set. The planner chose to fix it in the API: a project's web API reads that project's roadmap whatever all says, so every web caller (roadmap page, change set page's depends-on) is scoped, while the CLI's rhiza roadmap --all keeps widening across the root. The fix serves m37-roadmap ('The roadmap shows one project's milestones'). No question was raised.

## title

Plan of the bug fix: web roadmap page scoped to its project

## Links

### In

- derived-from [web-api-reads-own-project](../decision/web-api-reads-own-project.md)

## History

- `01M79S4M646SATZ2HW97GXK1MH` Fix The web roadmap page served under a project (/p/rhiza/roadmap) lists change sets and milestones from every registered project, not only that project's (owner, 2026-11-13: 'The roadmap tab doesn't filter between projects properly, it shows things from both projects'). loadRoadmap (internal/web/roadmap.go:105) reads /api/roadmap?all=true through the API handler in the request context; check whether that handler, or all=true, reaches across the root's projects instead of the scoped project. Wanted: each project's roadmap shows only its own change sets and milestones, with a test serving two projects' roadmaps and asserting no cross-project rows.
