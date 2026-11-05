# cpu-drivers-plan-interview

source · `01M6NEK8BY7P63JMRE5MZ27FBQ`

## quotes

```json
[]
```

## raw_sha256

6276853ef6393c1c77f461b89be762c0ceb0d2df9e85196dfbad88ceb67456b6

## session

a750dc81-b9fe-4c2e-904f-01a1aae0d697

## summary

Autonomous plan run for 'Fix Run drivers burn CPU while they wait'. The planner profiled 'rhiza change show' on this repo (2,539 log entries) with RHIZA_CPUPROFILE: about 8 s of CPU, ~5 s reading the log (store.ReadEntries: entry.Decode and canon.Canonicalize of every entry), ~1.5 s in changeset.FoldLog and ~2.3 s in changeNotes; driver poll paths (run_leaves, run_checks, run_review, run_build) read w.model() and w.changeSets(), and ship steps run 'rhiza ship check --settled' every 10 s. Live drivers were near idle while their runs sat in ship wait steps. No owner answered; the arbiter decided (assumed): the change set serves milestone-25, and model reads in drivers, 'rhiza change show', 'rhiza ship check --settled' and 'rhiza apply' go through the existing per-worktree SQLite index (internal/index) instead of reading and folding the whole log; an unchanged stat fingerprint means no model work; writes validate against the index's folded state plus their own entries; only 'rhiza check' and 'rhiza render' read the whole log; no second cache; an idle driver gets a CPU budget test. The planner wrote three leaves: idle driver polls, change show and ship check from the index, and apply from the index. Nothing is left open.

## title

Plan of the driver CPU fix, autonomous

## Links

### In

- derived-from [model-reads-go-through-the-index](../decision/model-reads-go-through-the-index.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
