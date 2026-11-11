# inbox-load-memory-hold-verified

evidence · `01M7635725YRD18NK96EB0Z7Y2` · actor [claude-code-e08303b5](../actor/claude-code-e08303b5.md)

## result

PASS

## statement

Verifier at b896e803: read the diff (inboxMachineHold/machineHoldItems in internal/api/inbox.go gate on healthLoadWarnAfter=30m from daemon state machine_hold, map memory-low reason to memory-low else load-high, detail names the reading, the limit and 'held for'; machine-wide, no tab to open; inbox app titles 'Machine load high'/'Memory low'). Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(inbox-load-hold|describe-golden)$': PASS; inbox-load-hold.txtar shows no item for a load hold since 29m, exactly one load-high item with '...over the limit of 3.00: held for 31m0s' at 31m, none for memory hold at 29m, one memory-low item at 31m, and items [] once a later poll's state has no hold. go test -trimpath ./internal/tui/inboxapp/...: ok (TestMachineHoldItemsRender renders both rows). TestScripts/inbox* pass; go test -trimpath ./internal/api -run 'Inbox|Notice|Machine|Hold' -v: all 17 tests PASS (incl. TestInboxNoticeNamesNoHandleOrID with load-high/memory-low cases); go vet and gofmt clean. The full internal/api package hit the 10m timeout on a loaded shared host in TestReadBudget (a timing budget that waits for a quiet host), unrelated to this diff.

## Links

### Out

- verifies [inbox-load-memory-hold](../capability/rhiza/orchestration/inbox/inbox-list/inbox-load-memory-hold.md)

## History

- `01M75T742VGKA5XMRRXV4W6ZKP` The daemon holds new runs while the machine is overloaded, by CPU load and by memory, instead of relying on a fixed rhiza.daemon.maxRuns that the orchestrator lowers and raises by hand. Owner, 2026-11-11: asked whether the daemon self-regulates on CPU load and memory (it does not: only disk via rhiza.daemon.minFree, the usage governor, and the fixed cap), then: "yes file it". Evidence: on the shared 4-core VM the orchestrator changed maxRuns between 3 and 2 by hand at least eight times on 2026-11-08 to 2026-11-10 (load 16-22 from three runs' Go compiles and review fan-outs, and builders' CPU busy loops on 2026-11-09). Wanted: before starting or resuming a run (ready work, request plan runs, schedules' runs, continuations), the daemon reads the 1-minute and 5-minute load average and available memory (Linux /proc/loadavg and MemAvailable in /proc/meminfo; on macOS sysctl vm.loadavg and vm_stat, or skip the check where unsupported) and starts nothing new while load per core is over a limit or available memory is under a floor. It never stops or pauses running work. Settings in rhiza.toml beside minFree, with sensible defaults, e.g. rhiza.daemon.maxLoad (per core, default about 3.0) and rhiza.daemon.minMemory (default about 2G); 0 or empty disables each. Use hysteresis so the daemon does not flap: hold above the limit, resume only when the 5-minute average is back under it. Surface it like the disk hold: rhiza daemon status shows the current load, memory and whether new runs are held and why; rhiza health and the inbox show a hold that has lasted longer than a threshold (e.g. 30 minutes); the daemon log notes each hold and release. maxRuns stays as the upper bound. Tests must not depend on the real machine's load: inject the reading (a function or a test env var) so a testscript can show held and released starts.
