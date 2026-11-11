# inbox-load-memory-hold

capability · `01M75TJ8FW01MJMA1HKD3TGNN0` · parent [inbox-list](../inbox-list.md)

## acceptance

An inbox testscript with a daemon state recording a load hold since 31 minutes before RHIZA_TEST_NOW shows one load-high item in rhiza inbox --json, none for a hold since 29 minutes, one memory-low item for a memory hold since 31 minutes, and none once a later poll records the hold released; go test -trimpath ./internal/tui/inboxapp/... covers the item's rendering

## statement

rhiza inbox lists one load-high or memory-low item while the daemon's last poll records a hold that has lasted 30 minutes or more, naming the reading, the limit and how long it has held, and the inbox app shows it

## Links

### Out

- decided-by [daemon-load-memory-hold-shape](../../../../../decision/daemon-load-memory-hold-shape.md)

### In

- verifies [inbox-load-memory-hold-verified](../../../../../evidence/inbox-load-memory-hold-verified.md)
- verifies [inbox-load-memory-hold-qa-1](../../../../../evidence/inbox-load-memory-hold-qa-1.md)
- verifies [inbox-load-memory-hold-verified-2](../../../../../evidence/inbox-load-memory-hold-verified-2.md)
- verifies [pr-442-the-daemon-holds-new-runs](../../../../../evidence/pr-442-the-daemon-holds-new-runs.md)

## History

- `01M75T742VGKA5XMRRXV4W6ZKP` The daemon holds new runs while the machine is overloaded, by CPU load and by memory, instead of relying on a fixed rhiza.daemon.maxRuns that the orchestrator lowers and raises by hand. Owner, 2026-11-11: asked whether the daemon self-regulates on CPU load and memory (it does not: only disk via rhiza.daemon.minFree, the usage governor, and the fixed cap), then: "yes file it". Evidence: on the shared 4-core VM the orchestrator changed maxRuns between 3 and 2 by hand at least eight times on 2026-11-08 to 2026-11-10 (load 16-22 from three runs' Go compiles and review fan-outs, and builders' CPU busy loops on 2026-11-09). Wanted: before starting or resuming a run (ready work, request plan runs, schedules' runs, continuations), the daemon reads the 1-minute and 5-minute load average and available memory (Linux /proc/loadavg and MemAvailable in /proc/meminfo; on macOS sysctl vm.loadavg and vm_stat, or skip the check where unsupported) and starts nothing new while load per core is over a limit or available memory is under a floor. It never stops or pauses running work. Settings in rhiza.toml beside minFree, with sensible defaults, e.g. rhiza.daemon.maxLoad (per core, default about 3.0) and rhiza.daemon.minMemory (default about 2G); 0 or empty disables each. Use hysteresis so the daemon does not flap: hold above the limit, resume only when the 5-minute average is back under it. Surface it like the disk hold: rhiza daemon status shows the current load, memory and whether new runs are held and why; rhiza health and the inbox show a hold that has lasted longer than a threshold (e.g. 30 minutes); the daemon log notes each hold and release. maxRuns stays as the upper bound. Tests must not depend on the real machine's load: inject the reading (a function or a test env var) so a testscript can show held and released starts.
