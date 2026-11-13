# daemon-refreshes-guide-at-start-verified

evidence · `01M7A05Z4M9JK249358261GB8X` · actor [claude-code-86fd20ca](../actor/claude-code-86fd20ca.md)

## result

PASS

## statement

At bc2c25ced re-ran the Checks block and acceptance: go test -trimpath -count=1 ./internal/api/... -run DaemonGuideRefresh (TestDaemonGuideRefresh 3 subtests and TestDaemonGuideRefreshAtStart PASS, exit 0); go test -trimpath ./internal/api/... -run 'DaemonReexec|DaemonLoop' exit 0; go vet ./internal/api/ exit 0; gofmt -l clean; rhiza ship generate exit 0; rhiza render exit 0. Diff (9a869996a, bc2c25ced) adds refreshGuides in internal/api/daemon_guide.go, called in runDaemon before the first poll; tests show the older home and project guides (missing skills/rhiza-builder) current at the first poll (ctx.Reap) after a restart and after a fake re-exec into a newer version, each scope logged; a home with no guide gets none; a foreign rhiza folder is untouched; current guides are not rewritten; a git-committed project guide is not rewritten, git status --porcelain stays empty, and the scope is a failed refresh whose Err and Next name rhiza guide install --agent claude --project. Only git runs (fake exec otherwise); no Rex, agents or gh. guideScopes is left for the sibling leaf daemon-holds-steps-on-stale-guide.

## Links

### Out

- verifies [daemon-refreshes-guide-at-start](../capability/rhiza/orchestration/daemon/daemon-refreshes-guide-at-start.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
