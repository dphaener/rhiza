# inbox-load-memory-hold-verified-2

evidence · `01M764QJS02ZEKPR6HXZHP5Y4K` · actor [claude-code-43db0f72](../actor/claude-code-43db0f72.md)

## result

PASS

## statement

At 7166e92d (fixer commit adding machineHoldItems to daemonNotify) I read the diff and ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-load-hold|health-load-hold|daemon-load-hold)$': ok. inbox-load-hold.txtar writes daemon state fixtures and checks rhiza inbox --json: a load hold since 29 minutes gives no item, since 31 minutes gives exactly one load-high item naming the reading, limit and 'held for 31m0s'; memory hold since 29 minutes gives none, since 31 minutes gives one memory-low item; a later poll state with no hold gives empty items. go test -trimpath ./internal/tui/inboxapp/...: ok (renders load-high 'Machine load high' and memory-low 'Memory low' rows). go test -trimpath ./internal/api -run TestDaemonNotifyMachineHold: ok, a 31-minute hold of either kind notifies once over two polls with the expected message. The api Inbox subset's tests all passed; the package-level root-state guard flagged ~/code/.rhiza-root/usage.live.json, which live Claude Code statuslines rewrite, so this change did not cause it.

## Links

### Out

- verifies [inbox-load-memory-hold](../capability/rhiza/orchestration/inbox/inbox-list/inbox-load-memory-hold.md)

## History

- `01M75T742VGKA5XMRRXV4W6ZKP` The daemon holds new runs while the machine is overloaded, by CPU load and by memory, instead of relying on a fixed rhiza.daemon.maxRuns that the orchestrator lowers and raises by hand. Owner, 2026-11-11: asked whether the daemon self-regulates on CPU load and memory (it does not: only disk via rhiza.daemon.minFree, the usage governor, and the fixed cap), then: "yes file it". Evidence: on the shared 4-core VM the orchestrator changed maxRuns between 3 and 2 by hand at least eight times on 2026-11-08 to 2026-11-10 (load 16-22 from three runs' Go compiles and review fan-outs, and builders' CPU busy loops on 2026-11-09). Wanted: before starting or resuming a run (ready work, request plan runs, schedules' runs, continuations), the daemon reads the 1-minute and 5-minute load average and available memory (Linux /proc/loadavg and MemAvailable in /proc/meminfo; on macOS sysctl vm.loadavg and vm_stat, or skip the check where unsupported) and starts nothing new while load per core is over a limit or available memory is under a floor. It never stops or pauses running work. Settings in rhiza.toml beside minFree, with sensible defaults, e.g. rhiza.daemon.maxLoad (per core, default about 3.0) and rhiza.daemon.minMemory (default about 2G); 0 or empty disables each. Use hysteresis so the daemon does not flap: hold above the limit, resume only when the 5-minute average is back under it. Surface it like the disk hold: rhiza daemon status shows the current load, memory and whether new runs are held and why; rhiza health and the inbox show a hold that has lasted longer than a threshold (e.g. 30 minutes); the daemon log notes each hold and release. maxRuns stays as the upper bound. Tests must not depend on the real machine's load: inject the reading (a function or a test env var) so a testscript can show held and released starts.
