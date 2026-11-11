# daemon-load-memory-hold-shape

decision · `01M75TGCCVXAFW35TVZT9JRAFK`

## reasoning

All four recommendations follow the change set's intent text, which the owner filed with "yes file it": it names maxLoad per core (about 3.0) and minMemory (about 2G) beside minFree, hysteresis that holds above the limit and resumes only when the 5-minute average is back under it, Linux /proc/loadavg and MemAvailable, macOS sysctl vm.loadavg or skipping the check where unsupported, and a health/inbox threshold of e.g. 30 minutes. m4-orchestrated-by-rhiza is the fitting intent: its statement is the daemon orchestrating with no orchestrator session, and the evidence is the orchestrator hand-tuning maxRuns eight times on 2026-11-08..10. Skipping memory on macOS avoids parsing vm_stat page counts for a check the VM (Linux) is what needs; the disk guard already skips off unix (internal/daemon/disk_other.go), so skipping where unsupported matches existing practice. A fixed 30-minute threshold keeps the settings surface small; it can become a setting later if needed. Memory needs no separate hysteresis window since MemAvailable is a point reading and the run-start cadence already damps flapping.

## statement

The change set serves @m4-orchestrated-by-rhiza. The daemon holds new runs when the 1-minute load per core is over rhiza.daemon.maxLoad (default 3.0) and releases only once the 5-minute load per core is back under it; it holds while MemAvailable is under rhiza.daemon.minMemory (default 2G) and releases once it is back at or over it; 0 or empty disables each. Linux reads /proc/loadavg and MemAvailable from /proc/meminfo; macOS reads load via sysctl vm.loadavg and skips the memory check; other platforms skip both, and rhiza daemon status says the reading is unavailable. rhiza daemon status and the daemon log show every hold and release at once; rhiza health and the inbox report a load or memory hold once it has lasted 30 minutes, a fixed threshold with no setting.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-942fa777-1](../question/question-942fa777-1.md)
- derived-from [plan-daemon-load-memory-hold](../source/plan-daemon-load-memory-hold.md)

### In

- decided-by [daemon-load-memory-hold](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-load-memory-hold.md)
- decided-by [daemon-load-memory-settings](../capability/rhiza/orchestration/daemon/daemon-settings/daemon-load-memory-settings.md)
- decided-by [daemon-status-load-memory](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-load-memory.md)
- decided-by [health-load-memory-hold](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-report/health-load-memory-hold.md)
- decided-by [inbox-load-memory-hold](../capability/rhiza/orchestration/inbox/inbox-list/inbox-load-memory-hold.md)

## History

- `01M75T742VGKA5XMRRXV4W6ZKP` The daemon holds new runs while the machine is overloaded, by CPU load and by memory, instead of relying on a fixed rhiza.daemon.maxRuns that the orchestrator lowers and raises by hand. Owner, 2026-11-11: asked whether the daemon self-regulates on CPU load and memory (it does not: only disk via rhiza.daemon.minFree, the usage governor, and the fixed cap), then: "yes file it". Evidence: on the shared 4-core VM the orchestrator changed maxRuns between 3 and 2 by hand at least eight times on 2026-11-08 to 2026-11-10 (load 16-22 from three runs' Go compiles and review fan-outs, and builders' CPU busy loops on 2026-11-09). Wanted: before starting or resuming a run (ready work, request plan runs, schedules' runs, continuations), the daemon reads the 1-minute and 5-minute load average and available memory (Linux /proc/loadavg and MemAvailable in /proc/meminfo; on macOS sysctl vm.loadavg and vm_stat, or skip the check where unsupported) and starts nothing new while load per core is over a limit or available memory is under a floor. It never stops or pauses running work. Settings in rhiza.toml beside minFree, with sensible defaults, e.g. rhiza.daemon.maxLoad (per core, default about 3.0) and rhiza.daemon.minMemory (default about 2G); 0 or empty disables each. Use hysteresis so the daemon does not flap: hold above the limit, resume only when the 5-minute average is back under it. Surface it like the disk hold: rhiza daemon status shows the current load, memory and whether new runs are held and why; rhiza health and the inbox show a hold that has lasted longer than a threshold (e.g. 30 minutes); the daemon log notes each hold and release. maxRuns stays as the upper bound. Tests must not depend on the real machine's load: inject the reading (a function or a test env var) so a testscript can show held and released starts.
