# decision-guide-refresh-start-fail-health

decision · `01M79NCFAEPVAHFCPX692PAJ3H`

## reasoning

The change set's intent asks for three fixes and offers fix 1 as 'refresh, or refuse with an inbox item'. Refreshing is the better default. It is what the orchestrator did by hand to unstick the runs. autoUpdate already means the owner wants the daemon to follow main without being asked. Refusing alone would still have stalled every run until someone noticed. Refusing with an inbox item stays as the fallback for when a refresh fails, which covers both options in the intent. The refresh runs at every start, not only on re-exec, because a fresh daemon start with a stale guide fails the same way. Status already decides 'current' as version plus plugin (status.go pluginCurrent), so that is the test the refresh uses. Only installed scopes are refreshed, so the daemon never writes into a home or project that has not opted in. For fix 2, the claudecode adapter's hooks show SessionStart and idle but nothing for a rejected slash command. An unknown command is never submitted as a prompt, so the failure has to be caught ahead of time with a preflight check that the skill exists. The UserPromptSubmit grace is a backstop for other agents and other causes. The incident cost 16 attempts on a fault in the machine's setup, so this failure must not use up role attempts. Holding the run with an inbox item matches how other environment failures go to the owner. For fix 3, health already reports per part with ok, warn and bad, and an out-of-date guide is worth a look but does not stop rhiza from running. That makes it warn, as the intent says. A machine-level part fits because guides are installed per machine or user, like the usage hold. The new enum value is an addition and breaks no existing user of --json.

## statement

(1) The daemon refreshes guides every time it starts, including after it re-executes a new binary. It does this before it starts any agent step. For each installed guide that status finds (user scope, then project scope) whose version or plugin is not the binary's, it reinstalls that scope's guide for that agent, the same as rhiza guide install --agent <agent>. It never installs a guide where none was installed. A project-scope guide tracked by git in the project's checkout is not rewritten, so the main checkout stays clean; it is handled as a failed refresh. If a refresh fails, the daemon does not start agent steps whose start command names a skill missing from the installed guide. It opens one inbox item naming the failed scope and the command rhiza guide install --agent <agent>, and it starts those steps again once the guide is current. (2) Before it launches a session, the adapter checks that the skill named by the step's start command is installed for the agent. If the skill is missing, the step fails at once with the reason 'start command /<name> unknown to <agent>: skill not installed; run rhiza guide install --agent <agent>'. This is an environment failure, not the role's: it uses up no builder or verifier attempt, and the run is held with an inbox item until the guide is installed and the run is resumed. There is also a backstop. A session that goes idle after SessionStart with no UserPromptSubmit within a short grace fails the same way, with the reason 'start command was never submitted'. (3) rhiza health gets a new machine-level part, guide, listed after usage. It is warn while any installed guide is out of date or missing a skill that a workflow starts, with next: rhiza guide install --agent <agent>. Otherwise it is ok, and it is never bad.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-ce3d52b3-1](../question/question-ce3d52b3-1.md)
- derived-from [plan-fix-runs-started-agent](../source/plan-fix-runs-started-agent.md)

### In

- decided-by [daemon-refreshes-guide-at-start](../capability/rhiza/orchestration/daemon/daemon-refreshes-guide-at-start.md)
- decided-by [daemon-holds-steps-on-stale-guide](../capability/rhiza/orchestration/daemon/daemon-holds-steps-on-stale-guide.md)
- decided-by [session-start-skill-preflight](../capability/rhiza/runtime/rex-sessions/session-start-skill-preflight.md)
- decided-by [session-start-prompt-not-submitted](../capability/rhiza/runtime/rex-sessions/session-start-prompt-not-submitted.md)
- decided-by [health-guide-part](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-guide-part.md)

## History

- `01M79N2JBDRDZ6R6HZ5YMRE8EY` Fix Runs started agent sessions with a slash command whose skill was not installed, and those sessions sat idle, burning build attempts. After #521 (build role topics) merged at 03:06 on 2026-11-13, the root daemon auto-updated (rhiza.daemon.autoUpdate = main) to a binary whose builtin workflows start builders with /rhiza-builder and verifiers with /rhiza-verifier, but ~/.claude/skills/rhiza still held the older guide (rhiza-build, rhiza-plan, rhiza-review, rhiza-scrub; rhiza status already said 'Guide (user): out of date'). From 03:29 to 04:55, 16 sessions in runs 01M79EVFQQ, 01M79DZCC6 and 01M79DC5TX opened with Claude Code's 'Unknown command: /rhiza-builder' (12) or '/rhiza-verifier' (4) and never got their task; each cost a builder or verifier attempt (seen as 'builder fail'), and four were still idle at 04:55 (cb3bbb91, 3296ae2f, 3aff6de9, cac42031). Unstuck by the orchestrator: rhiza guide install --agent claude, then rhiza session message to each idle session with its original task. Wanted: (1) when the daemon re-executes a new binary, it refreshes the guide installs it finds out of date (or refuses to start agent steps whose skill is missing, with an inbox item naming rhiza guide install); (2) a session whose start command is unknown to the agent fails at once with that reason instead of sitting idle; (3) rhiza health warns while the installed guide is out of date.
