# session-start-skill-preflight-qa-2

evidence · `01M7A702T51K1KTNM45HE97ZA2` · actor [claude-code-c4fbc364](../actor/claude-code-c4fbc364.md)

## result

PASS

## statement

QA cycle 2, scratch root with fake rex/claude and the branch build (go build -trimpath to a scratch path): a build workflow whose builder prompt is /design:design-critique failed at once with 'start command /design:design-critique unknown to claude: skill not installed; run rhiza guide install --agent claude' and an inbox item naming rhiza run resume --retry while the plugin was absent; after writing plugins/cache/mkt/design/1.0.0/skills/design-critique/SKILL.md under HOME/.claude, rhiza run resume --foreground --retry launched the builder as attempt 1 (fake builder logged 'alpha 1') and the inbox emptied, so the cycle-1 plugin finding is fixed. go test -trimpath ./internal/adapter/... and TestScripts/run-start-skill-missing passed.

## Links

### Out

- verifies [session-start-skill-preflight](../capability/rhiza/runtime/rex-sessions/session-start-skill-preflight.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
