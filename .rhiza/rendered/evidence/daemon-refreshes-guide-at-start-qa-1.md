# daemon-refreshes-guide-at-start-qa-1

evidence · `01M7A3S4G6YJ9DA5HQD0PR22RT` · actor [claude-code-2501afd2](../actor/claude-code-2501afd2.md)

## result

PASS

## statement

QA cycle 1, scratch root with fake rex/claude and the branch's build (go build -trimpath to a scratch path): started rhiza daemon --permission-mode manual with a user guide at version 0.0.9 missing rhiza-verifier and an untracked project guide missing rhiza-builder; the daemon logged 'refreshed the user guide ... to <version>' and 'refreshed the project repo guide ...' before any run, both installs were back at the running version with the skills present. With the project guide committed (git add -f .claude) and stale, the daemon logged 'could not refresh the project repo guide ...: the guide is committed in the repo: refresh it with rhiza guide install --agent claude --project and commit it through a change set', left the files unchanged and git status --porcelain empty. A foreign rhiza folder (no plugin manifest) was not treated as a guide. go test -trimpath ./internal/api/... -run DaemonGuideRefresh passed.

## Links

### Out

- verifies [daemon-refreshes-guide-at-start](../capability/rhiza/orchestration/daemon/daemon-refreshes-guide-at-start.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
