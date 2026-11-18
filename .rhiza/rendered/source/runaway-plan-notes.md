# runaway-plan-notes

source · `01M7QN1K9VDSJ66N7D5TR41XR6`

## quotes

```json
[
  {
    "by": "owner",
    "text": "Yes that makes sense to me.",
    "where": "2026-11-17 orchestrator chat, relayed by the orchestrator, recorded in Decision m40-host-limits-as-holds"
  }
]
```

## raw_sha256

a570cede787bbc421d9d1a18ac7b3ee4503521142ddab6729ecdbbf1532637a2

## summary

No owner was present: the conductor planned change set 01M7NVMD under the autonomous rules from the landed milestone-40 Decision m40-host-limits-as-holds (owner: 'Yes that makes sense to me.') and a reading of the code. Findings: a session's process tree is rooted at the 'rhiza session run <session-id>' process (the agent, e.g. claude --session-id, is its child), so ps -A -o pid=,ppid=,rss=,args= finds it on Linux and macOS alike (rss in KiB on both; internal/quiettest already reads ps, not /proc, for that reason). The watchdog's stalled-step pass (internal/api/daemon_watchdog.go, watchdogStalled) already nudges the step's stall sessions (stallSessions: open spawned children counted) and closes panes and terminates the run; only the engine journals a stopped entry today. Watchdog settings live in watchdogNames (internal/api/watchdog_policy.go) with machine and project keys; size values parse with daemon.ParseSize like rhiza.daemon.minTmpFree. Choices made by the planner, no owner asked: three new watchdog keys (runawayProcs default 400, runawayMemory default 8G, runawayGrace default 5m, 0 turning a threshold off); a runaway session is nudged once with the reading, and if a poll at least the grace later still reads it over a threshold the watchdog closes the stall sessions' panes, terminates the run and journals the step stopped with the reading as its detail, without resuming or retrying it (a retry would rerun the same runaway), filing a bug when fileBugs is on; attended runs under autonomousOnly get a log line only, like stalls; tests never read the host's process table.

## title

Planning pass for the watchdog's runaway sessions (autonomous conduct run)

## Links

### In

- derived-from [runaway-tree-from-session-run](../decision/runaway-tree-from-session-run.md)
- derived-from [runaway-thresholds-defaults](../decision/runaway-thresholds-defaults.md)
- derived-from [runaway-fails-not-retries](../decision/runaway-fails-not-retries.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason
