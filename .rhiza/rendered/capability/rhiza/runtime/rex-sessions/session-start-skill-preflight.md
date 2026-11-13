# session-start-skill-preflight

capability · `01M79NEZA13RZT7QH1TCZ113ER` · parent [rex-sessions](../rex-sessions.md)

## acceptance

Adapter tests (go test -trimpath ./internal/adapter/...) cover the skill lookup for an installed, a missing and a non-slash prompt; a testscript cmd/rhiza/testdata/script/run-start-skill-missing.txtar (go test -trimpath ./cmd/rhiza -run TestScripts/run-start-skill-missing) runs a build with HOME holding a guide without rhiza-builder and shows no rex launch, the step failed with that reason, the attempt count unchanged, the run held with an inbox item naming rhiza guide install --agent claude, and after installing the guide and resuming, the builder launches as attempt 1.

## statement

Before an agent step launches a session whose prompt starts with a slash command, the agent's adapter checks that the command's skill is installed for that agent (in the user's or the worktree's skills folder); when it is not, nothing is launched and the step fails at once with the reason 'start command /<name> unknown to <agent>: skill not installed; run rhiza guide install --agent <agent>', as an environment failure that uses up no builder, verifier or other role attempt and holds the run with an inbox item naming that command until the guide is installed and the run is resumed.

## Links

### Out

- decided-by [decision-guide-refresh-start-fail-health](../../../../decision/decision-guide-refresh-start-fail-health.md)

### In

- verifies [session-start-skill-preflight-verified](../../../../evidence/session-start-skill-preflight-verified.md)
- verifies [session-start-skill-preflight-qa-1](../../../../evidence/session-start-skill-preflight-qa-1.md)
- verifies [session-start-skill-preflight-verified-2](../../../../evidence/session-start-skill-preflight-verified-2.md)
- verifies [session-start-skill-preflight-qa-2](../../../../evidence/session-start-skill-preflight-qa-2.md)
- verifies [pr-540-fix-runs-started-agent](../../../../evidence/pr-540-fix-runs-started-agent.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
