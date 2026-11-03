# m24-run-drivers-interview

source · `01M6FEJYGZZC4DKW61YDRW04HB`

## quotes

```json
[
  {
    "by": "owner",
    "text": "All these tabs that are open right now, whatever is watching the plan, that should also be swallowed up by the daemon. No need to see anything in a Rex tab that isn't doing agent work.",
    "where": "message during the interview"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Detached child of the daemon: the daemon spawns rhiza run start as a background process with its output in a per-run log file under .git/rhiza; a run still outlives a daemon restart; each agent step opens its own Rex tab"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "All go through the daemon: the roadmap app and the CLI ask the daemon to start the run, so every driver is headless and counted against maxRuns; rhiza run start --foreground stays for debugging"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Existing apps + log: add rhiza run log <run> and a key in the runs app that opens it in a Rex layer on demand"
  }
]
```

## raw_sha256

0574da2ac0b6b8bbb6c82ce5d7bf4efc3f462a47cdf72cd54fc7ceaf1ccf16b8

## session

841f57a3-6242-4b77-8c1f-01a19edced79

## summary

Mid-interview the owner added that the Rex tabs watching each plan or build run should move into the daemon too, so Rex shows only tabs doing agent work. Settled: a run's driver becomes a detached child process of the daemon with its output in a per-run log file, so a run still outlives a daemon restart; each agent step opens its own Rex tab; every run starts through the daemon, including the roadmap app's start keys and rhiza run start, with a foreground mode kept for debugging; progress shows in the existing apps plus rhiza run log and a runs-app key that opens the log. This becomes a fourth change set of the milestone, planned now.

## title

Milestone 24 plan, part two: run drivers move into the daemon

## Links

### In

- derived-from [m24-only-agent-tabs](../decision/m24-only-agent-tabs.md)
- derived-from [answer-841f57a3-1](../decision/answer-841f57a3-1.md)
- derived-from [answer-841f57a3-2](../decision/answer-841f57a3-2.md)
- derived-from [answer-841f57a3-3](../decision/answer-841f57a3-3.md)
- derived-from [answer-841f57a3-4](../decision/answer-841f57a3-4.md)
- derived-from [answer-841f57a3-5](../decision/answer-841f57a3-5.md)
- derived-from [answer-841f57a3-6](../decision/answer-841f57a3-6.md)
- derived-from [answer-841f57a3-7](../decision/answer-841f57a3-7.md)
- derived-from [answer-841f57a3-8](../decision/answer-841f57a3-8.md)
- derived-from [answer-841f57a3-9](../decision/answer-841f57a3-9.md)
- derived-from [answer-841f57a3-10](../decision/answer-841f57a3-10.md)
- derived-from [answer-841f57a3-11](../decision/answer-841f57a3-11.md)
- derived-from [answer-841f57a3-12](../decision/answer-841f57a3-12.md)
- derived-from [answer-841f57a3-13](../decision/answer-841f57a3-13.md)

## History

- `01M6FDSQH550AN63HK4MK00T6V` Make Rex run as a background daemon on Linux and move this installation to a dedicated VM
