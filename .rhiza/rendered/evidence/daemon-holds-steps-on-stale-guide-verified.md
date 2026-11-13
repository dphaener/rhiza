# daemon-holds-steps-on-stale-guide-verified

evidence · `01M7A2ZMNXPDSPX3118753VC21` · actor [claude-code-fc0c2a01](../actor/claude-code-fc0c2a01.md)

## result

PASS

## statement

At 4d9ba18ce re-ran the commit's Checks: go test -trimpath ./internal/api/... -run DaemonGuideHold (TestDaemonGuideHold PASS, exit 0); go test -trimpath ./internal/api/... (exit 0, 864s with -timeout 40m: the default 10m timed out in unrelated TestUsageHoldStopsEveryLaunch under host load 13 on 4 cores); go test -trimpath ./internal/daemon/... (exit 0); rhiza ship generate and rhiza render (exit 0, tree unchanged); plus TestScripts describe-golden/inbox (exit 0). Read the diff: guideGate keeps the failed or tracked-skipped scopes from refreshGuides, guideSkipOf skips a change set whose workflow start command names a guide skill missing from every installed guide (SkipGuideStale, in daemonPoll1 and daemonDemand), state.GuideHolds feeds one guide-stale inbox item per scope naming the scope, the error (the committed-in-repo text for a tracked project guide) and rhiza guide install --agent claude (--project for a project scope), and holds drop once pluginCurrent. The test makes the user skills folder unwritable, shows no start, the skip, and exactly one item over 3 polls, then after guide install the next poll starts the target with no item. No stubs or TODOs; changes stay in the leaf's scope.

## Links

### Out

- verifies [daemon-holds-steps-on-stale-guide](../capability/rhiza/orchestration/daemon/daemon-holds-steps-on-stale-guide.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
