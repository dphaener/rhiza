# health-guide-part-qa-1

evidence · `01M7A3S4MAGJ3HG9KRYM9ERPT8` · actor [claude-code-2501afd2](../actor/claude-code-2501afd2.md)

## result

PASS

## statement

QA cycle 1, scratch HOME and project with the branch's build: rhiza health --json gave guide ok '0 installed guide(s) current' with no guide, ok '1 installed guide(s) current' after guide install; warn naming the user guide and next rhiza guide install --agent claude with skills/rhiza-builder removed and with version 0.0.9 ('version 0.0.9 is not <version>'); warn with next rhiza guide install --agent claude --project for a stale project guide; a foreign rhiza folder counted as no guide (ok). Human output lists guide right after usage; with only warns, rhiza health exited 0. go test -trimpath ./internal/api/... -run HealthGuide and TestScripts/describe-golden passed.

## Links

### Out

- verifies [health-guide-part](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-guide-part.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
