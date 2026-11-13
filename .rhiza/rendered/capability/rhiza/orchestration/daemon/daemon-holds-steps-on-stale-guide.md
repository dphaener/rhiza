# daemon-holds-steps-on-stale-guide

capability · `01M79NEZ70EHGQE9VERBQ04GFJ` · parent [daemon](../daemon.md)

## acceptance

Go tests in internal/api (go test -trimpath ./internal/api/... -run DaemonGuideHold) make the refresh fail (an unwritable skills folder) and show: an accepted change set's build step is not started, exactly one inbox item names the scope and rhiza guide install --agent claude across several polls, and after the guide is installed the next poll starts the step and the item is gone.

## statement

When the daemon's guide refresh fails for a scope, or it skips a project scope because the guide there is tracked by git, the daemon starts no agent step whose start command names a skill missing from that installed guide, opens one inbox item naming the scope, the error (or that the guide is committed in the repo and must be refreshed and committed through a change set) and the command rhiza guide install --agent <agent> (with --project for a project scope), and starts those steps on a later poll once the guide is current, clearing the item.

## Links

### Out

- decided-by [decision-guide-refresh-start-fail-health](../../../../decision/decision-guide-refresh-start-fail-health.md)

### In

- verifies [daemon-holds-steps-on-stale-guide-verified](../../../../evidence/daemon-holds-steps-on-stale-guide-verified.md)
- verifies [daemon-holds-steps-on-stale-guide-qa-1](../../../../evidence/daemon-holds-steps-on-stale-guide-qa-1.md)
- verifies [daemon-holds-steps-on-stale-guide-verified-2](../../../../evidence/daemon-holds-steps-on-stale-guide-verified-2.md)
- verifies [daemon-holds-steps-on-stale-guide-qa-2](../../../../evidence/daemon-holds-steps-on-stale-guide-qa-2.md)
- verifies [daemon-holds-steps-on-stale-guide-verified-3](../../../../evidence/daemon-holds-steps-on-stale-guide-verified-3.md)
- verifies [pr-540-fix-runs-started-agent](../../../../evidence/pr-540-fix-runs-started-agent.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
