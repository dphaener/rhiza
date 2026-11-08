# milestone-29

intent · `01M6RN3786W5DAVS69HMZ28RD9`

## milestone

29

## statement

Rhiza manages and builds several projects on one machine from one root directory: projects are registered under the root, one root daemon builds them all under one run cap and one usage budget, the CLI, web UI and TUI apps show every project and each on its own, and this VM runs rhiza and nebari that way.

## Links

### Out

- decided-by [m29-one-root-daemon](../decision/m29-one-root-daemon.md)
- decided-by [m29-plan-all-now](../decision/m29-plan-all-now.md)

### In

- serves change set `01M6RN3M8EC8FEDEJ8X6YVCPZ4` (Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project)
- serves change set `01M6RN3RCGAZ9RNTWTZPKFM61M` (Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service)
- serves change set `01M6RN3VGZ7D1JK4FRPNX21WHX` (Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one)
- serves change set `01M6RN3ZBGYDW6KPBSES9A5JXS` (Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher)
- serves change set `01M6RN441P2Z865ZVKAS7X7GMG` (Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox)
- serves change set `01M6RN474B5EXMX8NYHE7F62RV` (This VM on one root: rhiza and nebari registered under one root run by one root daemon service, replacing rhiza-rhiza.service)
- serves change set `01M6VY7GRV7BAJ9MMFVAY9GPZ2` (Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.)
- serves change set `01M6W7G8YBXBX090SSZC3Y60H2` (Fix Review question handles collide across runs started together: reviewQuestionHandle (internal/api/run_review.go ~927) builds review-<rid[:8]>-r<round>-<finding>, but a run ID is a ULID whose first 8 characters are only its start time (~35 s granularity), so plan runs the daemon starts in the same poll share the prefix. Evidence: plan runs of change sets 01M6VY7G6HPJ (dispatch) and 01M6VY7GRV7B (scheduling) both added question review-01m6vy7e-r3-f1 (nodes 01M6VZP9EAZX and 01M6VZK56QGX); rhiza health warns 'in-flight: model: duplicate handle', and when both change sets land rhiza check will fail on main. Use the run ID's random part (or the full ID, or the change set) in the handle, and when a handle is taken pick the next free one; add a test with two runs whose IDs share their first 8 characters. Existing duplicates: retitle one of the two when it lands (not unstuck yet).)
- serves change set `01M6WGEPSYWV8R14VGKJX9GGJM` (Fix A Go build inside a run step fails when its $WORK temp folder vanishes mid-build, and the run ends: plan run 01M6W7G5J0M6 (change set 01M6W7G8YBXB) failed at step rebase on 2026-11-08 ~00:15Z with 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' for many std packages while ~5 agents built and tested on the 4-core VM (load 28). The daemon's go build sweep (gosweep, TempMaxAge 2h) logged nothing in that window, so something else removed a live ~/tmp/go-build* folder: find what (candidates: a test's TestMain or tmpsweep removing TMPDIR entries by prefix, an agent's cleanup, go's own cleanup under a shared TMPDIR) and stop it. Also: a step whose command fails with this signature should retry once instead of ending the run. Unstuck by rhiza run resume --retry 01M6W7G5J0M6.)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
