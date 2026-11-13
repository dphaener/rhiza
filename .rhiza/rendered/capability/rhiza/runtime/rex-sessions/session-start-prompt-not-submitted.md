# session-start-prompt-not-submitted

capability · `01M79NEZJ5H07ZV0CGASQ0W2MJ` · parent [rex-sessions](../rex-sessions.md)

## acceptance

Go tests (go test -trimpath ./internal/session/... ./internal/api/... -run StartNotSubmitted) feed hook reports to a run's session: SessionStart then idle past the grace with no UserPromptSubmit fails the session with that reason and holds the run without counting an attempt; SessionStart, UserPromptSubmit then idle does not; idle within the grace does not.

## statement

A run's session whose agent reports SessionStart and then goes idle with no UserPromptSubmit within a short grace fails at once with the reason 'start command was never submitted', handled as the same environment failure as a missing start skill (no role attempt used, run held with an inbox item), instead of sitting idle until the stall limit.

## Links

### Out

- decided-by [decision-guide-refresh-start-fail-health](../../../../decision/decision-guide-refresh-start-fail-health.md)

### In

- verifies [session-start-prompt-not-submitted-verified](../../../../evidence/session-start-prompt-not-submitted-verified.md)
- verifies [session-start-prompt-not-submitted-qa-1](../../../../evidence/session-start-prompt-not-submitted-qa-1.md)
- verifies [pr-540-fix-runs-started-agent](../../../../evidence/pr-540-fix-runs-started-agent.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
