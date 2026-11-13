# health-guide-part-verified

evidence · `01M79ZD0B8T43M379J3P0DT6F2` · actor [claude-code-85d83a20](../actor/claude-code-85d83a20.md)

## result

PASS

## statement

At 1cec42d08 re-ran every Checks: command, all exit 0: go test -trimpath ./internal/api/... -run HealthGuide (TestHealthGuide, TestHealthGuideAfterUsage PASS), go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden, go test -trimpath ./internal/api/... -run Health, go test -trimpath ./cmd/rhiza -run 'TestScripts/(health|daemon-version-in-health|watchdog-health)', rhiza ship generate (0 changed), rhiza render (no tracked diff); go vet and gofmt clean. Diff read: healthGuide in internal/api/health.go is placed right after healthUsage, checks the user scope and each covered project root via installedVersion/pluginCurrent, warns with detail per stale install and next 'rhiza guide install --agent claude' (plus --project for a project scope), is ok with none installed, never bad. Tests cover none (ok), current (ok), older version (warn+next), missing skills/rhiza-builder (warn), stale project scope (warn with --project), order after usage and exit 0. describe-golden and docs/reference/cli.md updated for the guide enum value. No stubs or TODOs; nothing outside the leaf changed besides the builder's actor record.

## Links

### Out

- verifies [health-guide-part](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-guide-part.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
