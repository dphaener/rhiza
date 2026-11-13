# session-start-prompt-not-submitted-verified

evidence · `01M7A0C59RVZ2WAH4DYXWWDX71` · actor [claude-code-67d314e5](../actor/claude-code-67d314e5.md)

## result

PASS

## statement

At 1dbf71699 re-ran every Checks: command (exit 0 each): go test -trimpath ./internal/session/... ./internal/api/... -run StartNotSubmitted (acceptance; TestStartNotSubmitted, TestStartNotSubmittedFailsSession, ...PromptedIsFine, ...WithinGrace pass), go test -trimpath ./internal/run/... -run StartNotSubmitted (TestLeavesStepStartNotSubmitted passes), go test -trimpath ./internal/adapter/..., go test -trimpath ./internal/api/ -run 'TestSessionHook|TestSessionRun|Lost', rhiza ship generate, rhiza render (no tree changes). Read the leaf's three commits (5e6306336, 6bdf97a51, 1dbf71699): claude-code SessionStart sets HookReport.Start; recordHook marks Started/Prompted; session.NeverSubmitted with a 3m StartGrace; AgentStatus ends the session with 'start command was never submitted' as Environment; pollLeaf journals leaf-void after closing the pane and fails the run, so no leaf-attempt is counted and a retry restarts the builder at attempt 1. The api tests drive real hook payloads on a moved clock: SessionStart then idle past the grace fails with the reason; SessionStart, UserPromptSubmit, Stop past the grace stays live; idle within the grace stays live. Fakes only; no TODOs or stubs; changes stay within the leaf.

## Links

### Out

- verifies [session-start-prompt-not-submitted](../capability/rhiza/runtime/rex-sessions/session-start-prompt-not-submitted.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
