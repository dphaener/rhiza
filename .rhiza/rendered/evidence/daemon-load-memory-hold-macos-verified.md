# daemon-load-memory-hold-macos-verified

evidence · `01M7648GPC4HZD3169N595HJC0` · actor [claude-code-ead74a1d](../actor/claude-code-ead74a1d.md)

## result

PASS

## statement

Verifier at c57cac64 (review fix for the confirmed macOS finding). Read the diff: ReadMachine now takes load from readLoadavg, /proc/loadavg on !darwin (load_proc.go) and unix.SysctlRaw("vm.loadavg") on darwin (load_darwin.go), parsed by ParseSysctlLoadavg (three uint32 fixed-point averages, padding, int64 fscale at offset 16, matching struct loadavg on 64-bit); memory stays /proc-only so macOS skips only the memory check, as the shape decision says; an unreadable sysctl leaves LoadOK false and never holds. go.mod moves golang.org/x/sys v0.47.0 from indirect to direct, same version. Ran: go test -trimpath ./internal/daemon/... ok (includes TestParseSysctlLoadavg: 2.5/1.25 at scale 2048, short and zero-scale inputs fail); GOOS=darwin GOARCH=arm64 go vet ./internal/daemon and go build ./cmd/rhiza ok; gofmt -l clean; go test -trimpath ./cmd/rhiza -run TestScripts/daemon-load: daemon-load-hold and daemon-load-settings PASS. The darwin path itself was not run on a Mac.

## Links

### Out

- verifies [daemon-load-memory-hold](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-load-memory-hold.md)

## History

- `01M75T742VGKA5XMRRXV4W6ZKP` The daemon holds new runs while the machine is overloaded, by CPU load and by memory, instead of relying on a fixed rhiza.daemon.maxRuns that the orchestrator lowers and raises by hand. Owner, 2026-11-11: asked whether the daemon self-regulates on CPU load and memory (it does not: only disk via rhiza.daemon.minFree, the usage governor, and the fixed cap), then: "yes file it". Evidence: on the shared 4-core VM the orchestrator changed maxRuns between 3 and 2 by hand at least eight times on 2026-11-08 to 2026-11-10 (load 16-22 from three runs' Go compiles and review fan-outs, and builders' CPU busy loops on 2026-11-09). Wanted: before starting or resuming a run (ready work, request plan runs, schedules' runs, continuations), the daemon reads the 1-minute and 5-minute load average and available memory (Linux /proc/loadavg and MemAvailable in /proc/meminfo; on macOS sysctl vm.loadavg and vm_stat, or skip the check where unsupported) and starts nothing new while load per core is over a limit or available memory is under a floor. It never stops or pauses running work. Settings in rhiza.toml beside minFree, with sensible defaults, e.g. rhiza.daemon.maxLoad (per core, default about 3.0) and rhiza.daemon.minMemory (default about 2G); 0 or empty disables each. Use hysteresis so the daemon does not flap: hold above the limit, resume only when the 5-minute average is back under it. Surface it like the disk hold: rhiza daemon status shows the current load, memory and whether new runs are held and why; rhiza health and the inbox show a hold that has lasted longer than a threshold (e.g. 30 minutes); the daemon log notes each hold and release. maxRuns stays as the upper bound. Tests must not depend on the real machine's load: inject the reading (a function or a test env var) so a testscript can show held and released starts.
