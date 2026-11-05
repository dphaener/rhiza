# question-a750dc81-1

question · `01M6NEAK07BMQRE64NCYGHYF5Z`

## session

a750dc81-b9fe-4c2e-904f-01a1aae0d697

## statement

Change set 'Fix Run drivers burn CPU while they wait' serves no intent yet. Which intent should it serve? Earlier fix change sets (e.g. Fix Test temp folders leak) served milestone-25, the backlog intent.

## Links

### In

- answers [cpu-drivers-fix-serves-milestone-25](../decision/cpu-drivers-fix-serves-milestone-25.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
