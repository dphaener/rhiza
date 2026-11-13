# m37-roadmap

intent · `01M79NR64KCQDWD3N7HDVVBQ5R` · parent [milestone-37](../milestone-37.md)

## statement

The roadmap shows one project's milestones, their done points and change sets with short titles, state and progress, and every link on it opens a page that exists.

## Links

### In

- serves change set `01M79NMNGDCNJPAWRP791PRF7G` (Web roadmap and change set page: milestones, their done points and change sets with short titles, state and progress, scoped to the project; a change set page with its leaves, acceptance, Evidence, decisions, PR and runs)
- serves change set `01M79S4M646SATZ2HW97GXK1MH` (Fix The web roadmap page served under a project (/p/rhiza/roadmap) lists change sets and milestones from every registered project, not only that project's (owner, 2026-11-13: 'The roadmap tab doesn't filter between projects properly, it shows things from both projects'). loadRoadmap (internal/web/roadmap.go:105) reads /api/roadmap?all=true through the API handler in the request context; check whether that handler, or all=true, reaches across the root's projects instead of the scoped project. Wanted: each project's roadmap shows only its own change sets and milestones, with a test serving two projects' roadmaps and asserting no cross-project rows.)
- serves change set `01M79T7NW1B0KKXVVSX4AHZNHF` (Fix The web roadmap shows each change set's whole intent sentence as its title, so rows run off and say little (owner, 2026-11-13: 'The titles seem to be a description and just run off. It's not informative at all.'). roadmap.go:140 sets Row.Title to c.Intent; change sets carry a short title (rhiza change new --title, at most 30 characters, defaulting to the intent's first clause). Wanted: rows show the change set's title, with the intent as the secondary line clamped to two lines, on the roadmap and wherever else the web UI lists change sets (overview, runs list), with a test.)

## History

- `01M79NJ83BMM2XSVX2CMAWQXGG` Web UI redesign. Owner, 2026-11-13: the web UI is mostly useless: the roadmap mixes projects and its rows 404 and run off, run detail lists only Session 1, Session 2, the rendered model is used nowhere, the pages don't match their designs, spacing is off, and the header's project list needs a real picker. Redesign it screen by screen, phone-first, around 'does anything need me?', with a few safe actions.
