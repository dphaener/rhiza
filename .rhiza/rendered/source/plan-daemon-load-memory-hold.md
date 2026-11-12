# plan-daemon-load-memory-hold

source · `01M75THFY0GAA5DT6JEB3Z13SF`

## quotes

```json
[]
```

## raw_sha256

ab222623a70d3a613ef17deb54c5017e1292e425058a2fbd6caaf2577a5a3c44

## session

942fa777-a2c8-422f-8fd6-01a1cba4265c

## summary

Autonomous plan run for change set 01M75T742VGKA5XMRRXV4W6ZKP. No owner was present. The planner read the change set's intent (filed by the owner on 2026-11-11), the disk-low guard it should mirror (daemon-disk-guard, inbox-disk-low, healthDisk, RHIZA_TEST_FREE_BYTES) and the daemon settings, status and health capabilities. It raised four choices as one question for the arbiter: which intent the change set serves, the hold and release rule with defaults, how load and memory are read per platform, and when a hold reaches health and the inbox. The arbiter took every recommendation (Decision daemon-load-memory-hold-shape): serve m4-orchestrated-by-rhiza; hold when the 1-minute load per core is over rhiza.daemon.maxLoad (default 3.0) and release once the 5-minute load per core is back under it; hold while MemAvailable is under rhiza.daemon.minMemory (default 2G); 0 or empty disables each; Linux reads /proc, macOS reads load only, other platforms skip and say the reading is unavailable; daemon status and the log show every hold and release at once, health and the inbox once a hold has lasted 30 minutes, a fixed threshold. The planner then split the work into four leaves: the start hold itself, the two settings, the daemon status and log reporting, and the 30-minute health and inbox report.

## title

Plan: the daemon holds new runs on CPU load and memory

## Links

### In

- derived-from [daemon-load-memory-hold-shape](../decision/daemon-load-memory-hold-shape.md)

## History

- `01M75T742VGKA5XMRRXV4W6ZKP` The daemon holds new runs while the machine is overloaded, by CPU load and by memory, instead of relying on a fixed rhiza.daemon.maxRuns that the orchestrator lowers and raises by hand. Owner, 2026-11-11: asked whether the daemon self-regulates on CPU load and memory (it does not: only disk via rhiza.daemon.minFree, the usage governor, and the fixed cap), then: "yes file it". Evidence: on the shared 4-core VM the orchestrator changed maxRuns between 3 and 2 by hand at least eight times on 2026-11-08 to 2026-11-10 (load 16-22 from three runs' Go compiles and review fan-outs, and builders' CPU busy loops on 2026-11-09). Wanted: before starting or resuming a run (ready work, request plan runs, schedules' runs, continuations), the daemon reads the 1-minute and 5-minute load average and available memory (Linux /proc/loadavg and MemAvailable in /proc/meminfo; on macOS sysctl vm.loadavg and vm_stat, or skip the check where unsupported) and starts nothing new while load per core is over a limit or available memory is under a floor. It never stops or pauses running work. Settings in rhiza.toml beside minFree, with sensible defaults, e.g. rhiza.daemon.maxLoad (per core, default about 3.0) and rhiza.daemon.minMemory (default about 2G); 0 or empty disables each. Use hysteresis so the daemon does not flap: hold above the limit, resume only when the 5-minute average is back under it. Surface it like the disk hold: rhiza daemon status shows the current load, memory and whether new runs are held and why; rhiza health and the inbox show a hold that has lasted longer than a threshold (e.g. 30 minutes); the daemon log notes each hold and release. maxRuns stays as the upper bound. Tests must not depend on the real machine's load: inject the reading (a function or a test env var) so a testscript can show held and released starts.
