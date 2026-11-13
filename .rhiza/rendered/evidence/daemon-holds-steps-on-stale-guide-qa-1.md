# daemon-holds-steps-on-stale-guide-qa-1

evidence · `01M7A3S4JJ48QAMBHW516QEHZW` · actor [claude-code-2501afd2](../actor/claude-code-2501afd2.md)

## result

PASS

## statement

QA cycle 1, scratch root: an accepted change set, the user guide missing rhiza-builder and made read-only; rhiza daemon --autonomous logged 'could not refresh the user guide ...: permission denied', started no run (runs folder count unchanged), daemon status showed the change set skipped with reason guide-stale ('workflow ship-change starts /rhiza-builder, missing from the installed guide', the error, rhiza guide install --agent claude), and rhiza inbox held one guide-stale item naming the scope, the error and the command. After chmod u+w, the next daemon start refreshed the guide, started the ship-change run (builder launched as attempt 1) and the guide-stale item was gone. A committed stale project guide (no user guide) held the step with an item naming rhiza guide install --agent claude --project and the commit-through-a-change-set text. go test -trimpath ./internal/api/... -run DaemonGuideHold passed.

## Links

### Out

- verifies [daemon-holds-steps-on-stale-guide](../capability/rhiza/orchestration/daemon/daemon-holds-steps-on-stale-guide.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
