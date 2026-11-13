# daemon-refreshes-guide-at-start

capability · `01M79NEZ469MM69ZTQ6F8D0EK0` · parent [daemon](../daemon.md)

## acceptance

Go tests in internal/api (go test -trimpath ./internal/api/... -run DaemonGuideRefresh) start the daemon with a fake re-exec over a home and a project holding an older guide (missing skills/rhiza-builder), and show: both installs are rewritten to the running version before the first poll starts any session; a home with no guide gets none; a foreign rhiza folder is left untouched; a current guide is not rewritten; and a project whose older guide is committed to its git repo is not rewritten, git status --porcelain in its checkout stays empty, and the scope is reported as a failed refresh naming rhiza guide install --agent claude --project.

## statement

Each time the daemon starts, after a re-exec into a new binary as well as after a restart, and before it starts any agent step, it reinstalls every guide that status finds installed (user scope, then each project's scope) and not current for the running binary (version or plugin files differ), the same as rhiza guide install --agent <agent> for that scope; it never installs where no guide was installed and never overwrites a folder guide install did not write, and it logs each scope it refreshed. A project-scope guide with any file tracked by git in that project's checkout (git ls-files .claude/skills/rhiza is not empty, as a rhiza guide install --project committed to the repo is) is never rewritten: the daemon leaves the checkout untouched and treats that scope as a failed refresh (daemon-holds-steps-on-stale-guide), whose inbox item says to run rhiza guide install --agent <agent> --project and commit the result through a change set, so the main checkout stays clean for close-out's pull.

## Links

### Out

- decided-by [decision-guide-refresh-start-fail-health](../../../../decision/decision-guide-refresh-start-fail-health.md)

### In

- verifies [daemon-refreshes-guide-at-start-verified](../../../../evidence/daemon-refreshes-guide-at-start-verified.md)
- verifies [daemon-refreshes-guide-at-start-qa-1](../../../../evidence/daemon-refreshes-guide-at-start-qa-1.md)
- verifies [daemon-refreshes-guide-at-start-qa-2](../../../../evidence/daemon-refreshes-guide-at-start-qa-2.md)
- verifies [pr-540-fix-runs-started-agent](../../../../evidence/pr-540-fix-runs-started-agent.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
