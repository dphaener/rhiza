# session-start-skill-preflight-qa-1

evidence · `01M7A3S4PGVF71MQE88XENYQNY` · actor [claude-code-2501afd2](../actor/claude-code-2501afd2.md)

## result

PASS

## statement

QA cycle 1, scratch root with fake rex/claude: rhiza run start --foreground --workflow build (builder prompt /rhiza-builder) with HOME's guide lacking rhiza-builder failed at once with 'start command /rhiza-builder unknown to claude: skill not installed; run rhiza guide install --agent claude', no rex call, no session; inbox named the command and rhiza run resume --retry. Two further resumes failed the same way with no attempt used and one inbox item; after rhiza guide install --agent claude, resume --retry launched the builder as attempt 1 (fake builder logged 'beta 1') and cleared the inbox. A non-slash verifier prompt launched normally. go test -trimpath ./internal/adapter/... and TestScripts/run-start-skill-missing passed. One defect recorded as a finding: a marketplace-plugin command is falsely refused.

## Links

### Out

- verifies [session-start-skill-preflight](../capability/rhiza/runtime/rex-sessions/session-start-skill-preflight.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
