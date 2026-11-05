# cpu-drivers-fix-serves-milestone-25

decision · `01M6NECNPGK0JC3C1T1KVV1Q0T`

## reasoning

milestone-25 is where every recent fix and polish change set has gone (Fix Test temp folders leak, Fix Plan runs fail at merge when main moves, Fix Agents in runs hit permission prompts, and the comparable performance fix 'Inbox and daemon notify are cheap'). Its statement covers bugs and polish the daemon builds on its own. m24-on-the-vm held earlier VM fixes such as 'Fix The daemon pins a full CPU core', but milestone 24 is past and newer fixes stopped linking there. milestone-26 is about the TUI apps, which this fix does not touch. So milestone-25 fits the change set's intent and the repo's current practice best.

## statement

Change set 'Fix Run drivers burn CPU while they wait' serves milestone-25, the backlog intent

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-a750dc81-1](../question/question-a750dc81-1.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
