# daemon-load-memory-hold

capability · `01M75TJ828KDYB1XVYQ99EQVBM` · parent [daemon-picks](../daemon-picks.md)

## acceptance

A new testscript daemon-load-hold.txtar injects the readings through test env vars honoured only with RHIZA_TEST_NOW set (like RHIZA_TEST_FREE_BYTES): with the 1-minute load per core over maxLoad, rhiza daemon --once starts and resumes nothing and daemon status --json shows skips with reason load-high; with the 1-minute load back under but the 5-minute still over, it stays held; with the 5-minute under, the same fixture starts its run; MemAvailable under minMemory holds with reason memory-low; maxLoad 0 and minMemory 0 each disable their check; cmd/rhiza's test Setup gives every script a default reading so no script reads the host's load or memory; go test -trimpath ./internal/daemon/... covers parsing /proc/loadavg and /proc/meminfo samples

## statement

Before starting or resuming any run (ready work, plan requests, schedules' runs, continuations) the daemon reads load and available memory and starts nothing new while held: a hold begins when the 1-minute load per core is over rhiza.daemon.maxLoad and ends only once the 5-minute load per core is back under it, or while MemAvailable is under rhiza.daemon.minMemory; each held candidate is skipped with reason load-high or memory-low naming the reading and the limit; running runs are never stopped or paused; a reading that cannot be taken (macOS memory, other platforms) never holds; maxRuns stays the upper bound

## Links

### Out

- decided-by [daemon-load-memory-hold-shape](../../../../../decision/daemon-load-memory-hold-shape.md)

### In

- verifies [daemon-load-memory-hold-verified](../../../../../evidence/daemon-load-memory-hold-verified.md)
- verifies [daemon-load-memory-hold-qa-1](../../../../../evidence/daemon-load-memory-hold-qa-1.md)
- verifies [daemon-load-memory-hold-macos-verified](../../../../../evidence/daemon-load-memory-hold-macos-verified.md)
- verifies [pr-442-the-daemon-holds-new-runs](../../../../../evidence/pr-442-the-daemon-holds-new-runs.md)

## History

- `01M75T742VGKA5XMRRXV4W6ZKP` The daemon holds new runs while the machine is overloaded, by CPU load and by memory, instead of relying on a fixed rhiza.daemon.maxRuns that the orchestrator lowers and raises by hand. Owner, 2026-11-11: asked whether the daemon self-regulates on CPU load and memory (it does not: only disk via rhiza.daemon.minFree, the usage governor, and the fixed cap), then: "yes file it". Evidence: on the shared 4-core VM the orchestrator changed maxRuns between 3 and 2 by hand at least eight times on 2026-11-08 to 2026-11-10 (load 16-22 from three runs' Go compiles and review fan-outs, and builders' CPU busy loops on 2026-11-09). Wanted: before starting or resuming a run (ready work, request plan runs, schedules' runs, continuations), the daemon reads the 1-minute and 5-minute load average and available memory (Linux /proc/loadavg and MemAvailable in /proc/meminfo; on macOS sysctl vm.loadavg and vm_stat, or skip the check where unsupported) and starts nothing new while load per core is over a limit or available memory is under a floor. It never stops or pauses running work. Settings in rhiza.toml beside minFree, with sensible defaults, e.g. rhiza.daemon.maxLoad (per core, default about 3.0) and rhiza.daemon.minMemory (default about 2G); 0 or empty disables each. Use hysteresis so the daemon does not flap: hold above the limit, resume only when the 5-minute average is back under it. Surface it like the disk hold: rhiza daemon status shows the current load, memory and whether new runs are held and why; rhiza health and the inbox show a hold that has lasted longer than a threshold (e.g. 30 minutes); the daemon log notes each hold and release. maxRuns stays as the upper bound. Tests must not depend on the real machine's load: inject the reading (a function or a test env var) so a testscript can show held and released starts.
