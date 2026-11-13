# session-start-skill-preflight-verified

evidence · `01M79Z3KTC8BDXYV95YJJZS3H5` · actor [claude-code-e300208f](../actor/claude-code-e300208f.md)

## result

PASS

## statement

At 0e1e57d1e re-ran the Checks: block: go test -trimpath ./internal/adapter/... exit 0; go test -trimpath ./cmd/rhiza -run TestScripts/run-start-skill-missing exit 0; the builtin-workflow scripts pattern (bugfix-workflow-autonomous|usage-live-record|daemon-milestone-request|plan-milestone-workflow|workflow-plan-autonomy|run-start-skill-missing) exit 0; rhiza ship generate exit 0; rhiza render exit 0. Read the diff (63532bb2e): claudecode StartSkill/SkillGlobs with TestStartSkill (slash and non-slash prompts) and TestSkillGlobs (installed, missing, plugin and worktree skills); api checkStartSkill fails before launch with 'start command /<name> unknown to claude: skill not installed; run rhiza guide install --agent claude' and writes no session record. run-start-skill-missing.txtar asserts no rex call, run failed with that reason, no attempt, no sessions, an inbox item naming rhiza guide install --agent claude and resume --retry, and after reinstalling and resuming the builder runs as attempt 1. No TODOs or stubs; fakes only.

## Links

### Out

- verifies [session-start-skill-preflight](../capability/rhiza/runtime/rex-sessions/session-start-skill-preflight.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
