# m40-interview

source · `01M7NVKQKM6MQKPJJEMGQM9X8M`

## quotes

```json
[
  {
    "by": "owner",
    "text": "the overall idea here is that the daemon should be able to execute any workflow from start to finish without intervention. sure there may end up being errors from time to time that require intervention but for the most part it should just work. intervention should be the exception. even for bring your own workflows. so having an orchestrator like this should not be necessary. if that gets us closer let's plan it",
    "where": "2026-11-17 orchestrator chat"
  },
  {
    "by": "owner",
    "text": "A systemd scope seems like overkill to me.",
    "where": "2026-11-17 orchestrator chat, answer to question 1 (host limits)"
  },
  {
    "by": "owner",
    "text": "Yes that makes sense to me.",
    "where": "2026-11-17 orchestrator chat, accepting the orchestrator's host-limits proposal (spawn holds, free-tmp hold, runaway sessions as stalls)"
  },
  {
    "by": "planner (picked by the owner)",
    "text": "Right step + findings (Recommended)",
    "where": "2026-11-17, question 2 (failed runs)"
  },
  {
    "by": "planner (picked by the owner)",
    "text": "Yes, auto setting (Recommended)",
    "where": "2026-11-17, question 3 (concurrency)"
  },
  {
    "by": "owner",
    "text": "yes, push wake-up in, alarms folded into stall detection",
    "where": "2026-11-17 orchestrator chat, answer to question 4 (conductor wake-up)"
  }
]
```

## raw_sha256

6babc4b7c9ef9e2ddafece99977dce4e4e47fa780315c2f584c90a1451e4926b

## session

616e5d51-e7f1-4c75-ba67-01a1e9dadc50

## summary

On 2026-11-17 the owner asked that the daemon run any workflow, built-in or bring-your-own, from start to finish with intervention the exception, so the orchestrator session beside it becomes unnecessary. The orchestrator's brief listed what it does by hand today: machine rules sent to conductors, a stall watch, tmp sweeps, setting maxRuns by hand, retrying failed runs, and handing notes to planners. The planner asked four questions, and the owner answered through the orchestrator's chat. (1) Host limits: no systemd scopes or cgroups, which are Linux-only, and no daemon test-slot count. Instead the daemon's load, memory and disk holds also gate every child session spawn and add a free-tmp hold, and the stall watchdog treats a runaway session (its process tree's process count or memory past a host threshold, read portably) like a stall. The /tmp projection leak (01M7M3JZFW) and the 2.2 GB RSS of ship open-pr stay bug fixes. (2) Failed runs restart at the step the workflow names (default: the failed step; a plan refused by review-check restarts at plan), with the failure's output and findings in the next attempt's brief, one retry by default, then the inbox. (3) rhiza.daemon.maxRuns=auto starts another run only while the pre-start checks have room; a number still pins it. (4) The push wake-up draft 01M7M2FB5B serves this milestone; the check-in alarms of 01M7M2FE1S fold into the daemon's stall detection, and that draft is withdrawn. The planner ruled on what was left unasked: the order against the web split, the done measure, and planning each change set's leaves in its own plan run.

## title

Milestone 40 interview: the daemon runs any workflow, with the owner

## Links

### In

- derived-from [m40-no-systemd-scopes](../decision/m40-no-systemd-scopes.md)
- derived-from [m40-host-limits-as-holds](../decision/m40-host-limits-as-holds.md)
- derived-from [m40-retry-right-step](../decision/m40-retry-right-step.md)
- derived-from [m40-maxruns-auto](../decision/m40-maxruns-auto.md)
- derived-from [m40-push-wakeup-alarms-folded](../decision/m40-push-wakeup-alarms-folded.md)
- derived-from [m40-before-web-split](../decision/m40-before-web-split.md)
- derived-from [m40-hands-off-measure](../decision/m40-hands-off-measure.md)
- derived-from [m40-leaves-per-change-set](../decision/m40-leaves-per-change-set.md)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md
