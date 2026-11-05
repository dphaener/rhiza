# run-driver-idle-poll-cheap-qa-1

evidence · `01M6NYQT8WM94ZK683BQATY5CW` · actor [claude-code-edef540b](../actor/claude-code-edef540b.md)

## result

PASS

## statement

QA at 35496535, hands-on with binaries built from this branch (rhiza-new) and from the merge-base (rhiza-old), each run in its own scratch git repo that holds a copy of this repo's real .rhiza (2,617 entries, 51 MB), with RHIZA_ROOT/CHANGE/RUN unset. Ran a real `rhiza run start --foreground` driver (workflow: a leaves step, interval 2s), with a fake rex (from build-tabs.txtar) and a fake claude that sleeps, on PATH. The builder session started and the driver idle-polled it. Process CPU from /proc/<pid>/stat: 0.01 s over 40 s of idle polls (~20 polls). After an entry was appended to an unrelated change set in the run worktree, the next 6 s used 0 s. Main's driver also idled at ~0 here, because a leaves step reads leaves once per step entry. So this harness confirms the new driver idles cheap but does not reproduce the original burn. The Go budget covers the poll cost: TestReadBudget driver poll idle median 2.745 ms CPU (budget 5 ms), and after an unrelated entry 5.723 ms (budget 50 ms), PASS. go test -trimpath ./internal/api (run_model_cache_test.go among them): ok.

## Links

### Out

- verifies [run-driver-idle-poll-cheap](../capability/rhiza/orchestration/daemon/run-driver-idle-poll-cheap.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
