# session-start-skill-preflight-verified-2

evidence · `01M7A5921082PQHYRKE0D22H43` · actor [claude-code-d76a7f59](../actor/claude-code-d76a7f59.md)

## result

PASS

## statement

Attempt 2 at 40562d180: re-ran every Checks: command. TestSkillGlobs with claudecode.go reverted to HEAD~1 fails (exit 1: marketplace plugin skill/command not found); go test -trimpath ./internal/adapter/... exit 0 (TestStartSkill covers slash, non-slash and unsafe prompts; TestSkillGlobs covers installed, missing, worktree and plugin-cache/marketplace skills); go test -trimpath ./cmd/rhiza -run TestScripts/run-start-skill-missing exit 0 (no rex launch, failed step naming 'start command /rhiza-builder unknown to claude: skill not installed; run rhiza guide install --agent claude', no attempt, inbox item naming guide install, resume launches alpha attempt 1); go test -trimpath ./internal/api/ -run 'Skill|Guide' exit 0; full go test -trimpath ./internal/api/ exit 1 by 600s package timeout in TestReadBudgetApply, as recorded and outside this change; rhiza ship generate and rhiza render exit 0 with no tree changes. Diff limited to SkillGlobs plugin search plus its test and rendered actors; no TODOs or stubs.

## Links

### Out

- verifies [session-start-skill-preflight](../capability/rhiza/runtime/rex-sessions/session-start-skill-preflight.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
