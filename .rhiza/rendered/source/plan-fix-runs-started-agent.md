# plan-fix-runs-started-agent

source · `01M79NE48ZK6V86SVG7AK5HR74`

## quotes

```json
[]
```

## raw_sha256

ce044171fef262fda36e57c40e9548a9c62a4caf91e8fd109ef69603e4843d6f

## session

ce3d52b3-b257-41d5-a151-01a1d351d997

## summary

Autonomous plan run for the bug where runs started builder and verifier sessions with /rhiza-builder and /rhiza-verifier after the daemon auto-updated, while the installed guide was the older one: 16 sessions opened with 'Unknown command', each cost an attempt, and four sat idle. The planner read the daemon re-exec (internal/daemon/reexec.go, daemonLoop and daemonStartState in internal/api/daemon.go), the guide-currency test in status (installedVersion, pluginCurrent in internal/api/status.go), guide.install, the health parts, the session launch and the usage-limit transcript watch (internal/api/session_run.go). One transcript showed the agent writing a system message 'Unknown command: /rhiza-builder. Did you mean /rhiza-build?' and no prompt submitted. There was no owner. The planner raised three choices (when the daemon refreshes guides, what happens to a run whose start command is unknown, and how health reports a stale guide), and the arbiter answered with Decision decision-guide-refresh-start-fail-health: refresh installed guides at every daemon start, before any agent step; if a refresh fails, hold the agent steps whose start skill is missing, with one inbox item, until the guide is current; check before launch that the start command's skill is installed, and fail at once with an environment failure that uses no role attempt and holds the run in the inbox; as a backstop, fail a session that goes idle after SessionStart without UserPromptSubmit within a short grace; add a machine-level health part, guide, after usage, that is warn while a guide is stale or missing a skill a workflow starts, and never bad. Leaves were written from that Decision under the daemon, rex-sessions and health-command capabilities.

## title

Fix Runs started agent: autonomous plan interview

## Links

### In

- derived-from [decision-guide-refresh-start-fail-health](../decision/decision-guide-refresh-start-fail-health.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
