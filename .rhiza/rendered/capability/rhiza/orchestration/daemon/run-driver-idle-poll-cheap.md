# run-driver-idle-poll-cheap

capability · `01M6NEKX4E0HFPHC9ZSZKPE7KX` · parent [daemon](../daemon.md)

## acceptance

go test -trimpath ./internal/api/... : a Go test builds a repo with a large log (100+ change sets, several in flight) and its index (IndexPath set), drives a run's leaves and checks steps through the stepRunner with live sessions, polls repeatedly with nothing changed and asserts, through the ctx.readLog counter and decode and fold counters, that those polls read the log zero times, decode nothing and fold nothing; it then appends an entry to an unrelated change set and asserts the next poll reads the log zero times and decodes and folds only that change set; it then adds an entry to the run's change set and a session done mark and asserts the next poll sees both. A CPU budget test (quiettest, median, like TestReadBudget) measures the process CPU time of idle driver polls on the 50,000-node budget model, and of a poll after one unrelated appended entry, and asserts each costs under 5 ms and 50 ms of CPU respectively

## statement

A run driver's poll (agent, leaves and checks steps, and their arbiter waits) reads the model through the worktree's index, and a poll where the log's stat fingerprint, the run journal and the polled session records are unchanged does no model work: no log read, no decode and no fold

## Links

### Out

- depends-on [index-catchup-refolds-incrementally](../../views/fast-reads/index-catchup-refolds-incrementally.md)
- decided-by [model-reads-go-through-the-index](../../../../decision/model-reads-go-through-the-index.md)

### In

- verifies [run-driver-idle-poll-cheap-verified](../../../../evidence/run-driver-idle-poll-cheap-verified.md)
- verifies [run-driver-idle-poll-cheap-qa-1](../../../../evidence/run-driver-idle-poll-cheap-qa-1.md)
- verifies [run-driver-idle-poll-cheap-verified-2](../../../../evidence/run-driver-idle-poll-cheap-verified-2.md)
- verifies [pr-304-fix-run-drivers-burn-cpu-while](../../../../evidence/pr-304-fix-run-drivers-burn-cpu-while.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
