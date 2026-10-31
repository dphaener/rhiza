# m10-orchestrator-picks

source · `01M6A0KMK61TNC56042NHQRJ4H`

## quotes

```json
[
  {
    "by": "owner",
    "text": "if you run out of work keep picking high value tasks. don't stop and wait for my input.",
    "where": "orchestrator session, 2026-10-31"
  },
  {
    "by": "owner",
    "text": "keep monitoring and don't stop. restart the monitor when it stops, that is your entire job. to keep this running.",
    "where": "orchestrator session, 2026-10-30"
  }
]
```

## raw_sha256

92f31b5efccd32a19a68e7647ef599f0d87592238c9242111071c27c12193740

## summary

With milestones 7 to 9 queued, the orchestrator planned milestone 10 under the owner's standing instruction to keep picking high-value work. It is the hand fixes the orchestrator still made while runs went unattended: a PR that conflicts after main moves sits in a wait step with no CI until a hand rhiza ship rebase (GitHub runs no CI on a conflicting PR, and ship check never reads mergeable); a flaky CI check fails the run, and the orchestrator reruns it and merges by hand; a daemon start that never journaled (disk full) blocks its change set until the daemon is restarted, because the unjournaled map never expires; the daemon keeps running an old binary after daemon code lands, until restarted by hand; and agent sessions sat on a Claude Code dialog for hours, found only by eye, although each session's hook status carries the time of its last tool call. Research agents traced each to its code: ship.go's pr view fields and checksVerdict, stepRunner.Holds dropping exit 75, daemonHeld's unjournaled map, ctx.Self and the flock lock (dropped on exec, so a re-exec is safe once no start is pending), and recordHook's Status.Hook.At.

## title

Milestone 10: runs heal themselves

## Links

### In

- derived-from [m10-picked-by-orchestrator](../decision/m10-picked-by-orchestrator.md)
- derived-from [wait-settles-conflict-rule](../decision/wait-settles-conflict-rule.md)
- derived-from [ci-rerun-once-per-head](../decision/ci-rerun-once-per-head.md)
- derived-from [unjournaled-block-expires](../decision/unjournaled-block-expires.md)
- derived-from [daemon-reexec-when-idle](../decision/daemon-reexec-when-idle.md)
- derived-from [stalled-after-20-minutes](../decision/stalled-after-20-minutes.md)

## History

- `01M6A0JENBMBYPYRN5195AGPKH` Plan milestone 10: runs heal themselves, picked by the orchestrator under the owner's standing instruction
