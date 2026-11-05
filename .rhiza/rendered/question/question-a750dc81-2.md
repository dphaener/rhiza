# question-a750dc81-2

question · `01M6NEFDS6ZHEKC2RFK7QRFPEE`

## session

a750dc81-b9fe-4c2e-904f-01a1aae0d697

## statement

Profiling 'rhiza change show' on this repo (2,539 log entries) shows ~5 s of CPU spent reading the log (each entry decoded and re-canonicalized) and ~1.5 s folding it, every time. Drivers run in-process polls plus a 'rhiza ship check --settled' subprocess every 10 s, and 'rhiza apply' pays the full read per call. How should model reads be made cheap?

## Links

### In

- answers [model-reads-go-through-the-index](../decision/model-reads-go-through-the-index.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
