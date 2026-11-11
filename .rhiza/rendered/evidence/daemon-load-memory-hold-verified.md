# daemon-load-memory-hold-verified

evidence · `01M763RFPQ7253TM67RSGP1SY7` · actor [claude-code-960bdc88](../actor/claude-code-960bdc88.md)

## result

PASS

## statement

Verified at ab33bb96 (claude-opus verifier, attempt 2). Read the leaf diff (0df6bd6f + ab33bb96): internal/daemon/load.go parses /proc/loadavg and /proc/meminfo, Evaluate begins a load hold on the 1-minute load per core over maxLoad, keeps it until the 5-minute load is under (also beneath a memory-low hold), holds on MemAvailable under minMemory, 0 disables each, unknown readings never hold; daemonPoll1 treats a machine hold like disk-low (no plan requests, resumes, starts) with skips load-high/memory-low and their detail; schedules and demand stop too; running runs untouched; RHIZA_TEST_LOAD1/LOAD5/MEM_AVAILABLE apply only with RHIZA_TEST_NOW, and cmd/rhiza Setup sets idle defaults for every script. Ran: go test -trimpath ./internal/daemon/... -run 'TestParseLoadavg|TestParseMemAvailable|TestEvaluate' -v (PASS); go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-load-hold|describe-golden)$' (ok) - the script covers 1-minute over holds with load-high, 1-minute under/5-minute over stays held, 5-minute under starts the run, memory-low holds, minMemory 0 and maxLoad 0 each disable; go test -trimpath ./internal/api -run 'Daemon|Usage|Budget|Help' - all tests PASS, only TestMain's root-state guard reported ~/code/.rhiza-root/usage.live.json changed, which the live root daemon rewrites every few seconds (mtime 2s old), unrelated to this change.

## Links

### Out

- verifies [daemon-load-memory-hold](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-load-memory-hold.md)

## History

- `01M75T742VGKA5XMRRXV4W6ZKP` The daemon holds new runs while the machine is overloaded, by CPU load and by memory, instead of relying on a fixed rhiza.daemon.maxRuns that the orchestrator lowers and raises by hand. Owner, 2026-11-11: asked whether the daemon self-regulates on CPU load and memory (it does not: only disk via rhiza.daemon.minFree, the usage governor, and the fixed cap), then: "yes file it". Evidence: on the shared 4-core VM the orchestrator changed maxRuns between 3 and 2 by hand at least eight times on 2026-11-08 to 2026-11-10 (load 16-22 from three runs' Go compiles and review fan-outs, and builders' CPU busy loops on 2026-11-09). Wanted: before starting or resuming a run (ready work, request plan runs, schedules' runs, continuations), the daemon reads the 1-minute and 5-minute load average and available memory (Linux /proc/loadavg and MemAvailable in /proc/meminfo; on macOS sysctl vm.loadavg and vm_stat, or skip the check where unsupported) and starts nothing new while load per core is over a limit or available memory is under a floor. It never stops or pauses running work. Settings in rhiza.toml beside minFree, with sensible defaults, e.g. rhiza.daemon.maxLoad (per core, default about 3.0) and rhiza.daemon.minMemory (default about 2G); 0 or empty disables each. Use hysteresis so the daemon does not flap: hold above the limit, resume only when the 5-minute average is back under it. Surface it like the disk hold: rhiza daemon status shows the current load, memory and whether new runs are held and why; rhiza health and the inbox show a hold that has lasted longer than a threshold (e.g. 30 minutes); the daemon log notes each hold and release. maxRuns stays as the upper bound. Tests must not depend on the real machine's load: inject the reading (a function or a test env var) so a testscript can show held and released starts.
