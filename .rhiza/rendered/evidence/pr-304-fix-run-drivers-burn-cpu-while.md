# pr-304-fix-run-drivers-burn-cpu-while

evidence · `01M6P9E5H2HYBC5T4T85A8SN8Y` · actor [dphaener](../actor/dphaener.md)

## head

95e5f29ed00ae26f0ea1f7c413a94862328fb3a7

## pr

https://github.com/dphaener/rhiza/pull/304

## result

1 check passed: check SUCCESS

## statement

Pull request #304 (https://github.com/dphaener/rhiza/pull/304) passed CI at 95e5f29

## status

verified

## Links

### Out

- verifies [run-driver-idle-poll-cheap](../capability/rhiza/orchestration/daemon/run-driver-idle-poll-cheap.md)
- verifies [change-show-ship-check-from-index](../capability/rhiza/views/fast-reads/change-show-ship-check-from-index.md)
- verifies [apply-validates-from-index](../capability/rhiza/views/fast-reads/apply-validates-from-index.md)
- verifies [index-catchup-refolds-incrementally](../capability/rhiza/views/fast-reads/index-catchup-refolds-incrementally.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
