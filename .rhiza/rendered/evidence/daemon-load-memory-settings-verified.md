# daemon-load-memory-settings-verified

evidence · `01M761KEB38E8NRCT9R11EZWH2` · actor [claude-code-becf3bb8](../actor/claude-code-becf3bb8.md)

## result

PASS

## statement

At 88cc562d: go test -trimpath ./cmd/rhiza -run TestScripts/(daemon-load-settings|describe-golden|daemon-settings|daemon-settings-scope|daemon-load-hold|daemon-load-status) passes (describe-golden matches the regenerated max_load/min_memory fields and the daemon.set key list). By hand with a scratch build in a fresh root + project: rhiza daemon settings --json lists rhiza.daemon.maxLoad default 3.0 and rhiza.daemon.minMemory default 2G (machine scope); daemon set maxLoad 2.5 and minMemory 4G write them to the root rhiza.toml [daemon]; daemon set maxLoad -- -1, maxLoad busy, minMemory lots and minMemory 4Q each exit 2 with a validation message and cmp shows rhiza.toml unchanged; maxLoad 0 and minMemory 0 are accepted. Diff read: DaemonRequest gains MaxLoad/MinMemory (--max-load/--min-memory), validated in daemonConfOf with Usagef, settingDefs map the flags. Scoped internal/api tests (TestDaemon|TestSettingDefsDescribed|Usage|Hold|Load) all PASS; the package's post-run guard flagged ~/code/.rhiza-root/usage.live.json changed, which is the live root daemon on the VM, not this diff.

## Links

### Out

- verifies [daemon-load-memory-settings](../capability/rhiza/orchestration/daemon/daemon-settings/daemon-load-memory-settings.md)

## History

- `01M75T742VGKA5XMRRXV4W6ZKP` The daemon holds new runs while the machine is overloaded, by CPU load and by memory, instead of relying on a fixed rhiza.daemon.maxRuns that the orchestrator lowers and raises by hand. Owner, 2026-11-11: asked whether the daemon self-regulates on CPU load and memory (it does not: only disk via rhiza.daemon.minFree, the usage governor, and the fixed cap), then: "yes file it". Evidence: on the shared 4-core VM the orchestrator changed maxRuns between 3 and 2 by hand at least eight times on 2026-11-08 to 2026-11-10 (load 16-22 from three runs' Go compiles and review fan-outs, and builders' CPU busy loops on 2026-11-09). Wanted: before starting or resuming a run (ready work, request plan runs, schedules' runs, continuations), the daemon reads the 1-minute and 5-minute load average and available memory (Linux /proc/loadavg and MemAvailable in /proc/meminfo; on macOS sysctl vm.loadavg and vm_stat, or skip the check where unsupported) and starts nothing new while load per core is over a limit or available memory is under a floor. It never stops or pauses running work. Settings in rhiza.toml beside minFree, with sensible defaults, e.g. rhiza.daemon.maxLoad (per core, default about 3.0) and rhiza.daemon.minMemory (default about 2G); 0 or empty disables each. Use hysteresis so the daemon does not flap: hold above the limit, resume only when the 5-minute average is back under it. Surface it like the disk hold: rhiza daemon status shows the current load, memory and whether new runs are held and why; rhiza health and the inbox show a hold that has lasted longer than a threshold (e.g. 30 minutes); the daemon log notes each hold and release. maxRuns stays as the upper bound. Tests must not depend on the real machine's load: inject the reading (a function or a test env var) so a testscript can show held and released starts.
