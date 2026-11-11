# daemon-load-memory-hold-qa-1

evidence · `01M763YSF1PFHQ2TH2TH5DYS45` · actor [claude-code-415ed582](../actor/claude-code-415ed582.md)

## result

PASS

## statement

QA at f255dc77 with a scratch build (go build -trimpath) in a fresh Rhiza root + project under a scratch HOME, fake rex on PATH. Real readings (no RHIZA_TEST_NOW): host /proc/loadavg 12.17 7.77 on 4 cores; with maxLoad 0.1 rhiza daemon --once started nothing, logged 'holding new runs: load-high: 1-minute load 3.04 per core, over the limit of 0.10', and skipped the accepted change set with reason load-high; a second poll kept the original since time with no new hold line. maxLoad 0 + minMemory 100G held with memory-low on the real MemAvailable (11.0 GiB). Injected readings with RHIZA_TEST_NOW + detach log: load1/load5 3.5/3.5 held; 1/3.5 stayed held ('5-minute load 3.50 ... not yet under'); 1/3.0 stayed held; 1/2.9 logged 'released the load-high hold' and started the run (detach log shows rhiza run start). maxLoad 0 and minMemory 0 disable. --max-load 0.01 on the command line held too; --max-load -2 exits 2.

## Links

### Out

- verifies [daemon-load-memory-hold](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-load-memory-hold.md)

## History

- `01M75T742VGKA5XMRRXV4W6ZKP` The daemon holds new runs while the machine is overloaded, by CPU load and by memory, instead of relying on a fixed rhiza.daemon.maxRuns that the orchestrator lowers and raises by hand. Owner, 2026-11-11: asked whether the daemon self-regulates on CPU load and memory (it does not: only disk via rhiza.daemon.minFree, the usage governor, and the fixed cap), then: "yes file it". Evidence: on the shared 4-core VM the orchestrator changed maxRuns between 3 and 2 by hand at least eight times on 2026-11-08 to 2026-11-10 (load 16-22 from three runs' Go compiles and review fan-outs, and builders' CPU busy loops on 2026-11-09). Wanted: before starting or resuming a run (ready work, request plan runs, schedules' runs, continuations), the daemon reads the 1-minute and 5-minute load average and available memory (Linux /proc/loadavg and MemAvailable in /proc/meminfo; on macOS sysctl vm.loadavg and vm_stat, or skip the check where unsupported) and starts nothing new while load per core is over a limit or available memory is under a floor. It never stops or pauses running work. Settings in rhiza.toml beside minFree, with sensible defaults, e.g. rhiza.daemon.maxLoad (per core, default about 3.0) and rhiza.daemon.minMemory (default about 2G); 0 or empty disables each. Use hysteresis so the daemon does not flap: hold above the limit, resume only when the 5-minute average is back under it. Surface it like the disk hold: rhiza daemon status shows the current load, memory and whether new runs are held and why; rhiza health and the inbox show a hold that has lasted longer than a threshold (e.g. 30 minutes); the daemon log notes each hold and release. maxRuns stays as the upper bound. Tests must not depend on the real machine's load: inject the reading (a function or a test env var) so a testscript can show held and released starts.
