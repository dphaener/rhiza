# health-guide-part

capability · `01M79NEZPB651MRDDMV3S2GYJZ` · parent [health-command](../health-command.md)

## acceptance

Go tests in internal/api (go test -trimpath ./internal/api/... -run HealthGuide) cover a current guide (ok), an older version (warn with the next command), a guide missing skills/rhiza-builder (warn), a project-scope stale guide (warn with --project) and no guide (ok), and show the part sits after usage and leaves the exit code 0; describe-golden is regenerated for the new part enum (go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden).

## statement

rhiza health reports a machine-level part, guide, listed after usage: warn while any installed guide (user scope, or a covered project's scope) is out of date for the running rhiza or lacks a skill a workflow's start command names, naming each such install with next rhiza guide install --agent <agent> (plus --project for a project scope); ok otherwise, including when no guide is installed; never bad.

## Links

### Out

- decided-by [decision-guide-refresh-start-fail-health](../../../../../../decision/decision-guide-refresh-start-fail-health.md)

### In

- verifies [health-guide-part-verified](../../../../../../evidence/health-guide-part-verified.md)
- verifies [health-guide-part-qa-1](../../../../../../evidence/health-guide-part-qa-1.md)
- verifies [health-guide-part-qa-2](../../../../../../evidence/health-guide-part-qa-2.md)
- verifies [pr-540-fix-runs-started-agent](../../../../../../evidence/pr-540-fix-runs-started-agent.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
