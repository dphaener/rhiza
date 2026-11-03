# daemon-idle-poll-no-fold-qa-1

evidence · `01M6HF857940T40BCKP95PKSXY` · actor [claude-code-797a9572](../actor/claude-code-797a9572.md)

## result

PASS

## statement

At dd9192ca I built the branch binary and the merge-base ebee728a binary outside PATH. I cloned this repo (1746 log entries) to a scratch dir, set rhiza.daemon.paused=true, selfUpdate=false, minFree=0 and an eventCommand, and ran each binary's real `rhiza daemon --permission-mode auto` in the foreground. I read the process's utime+stime from /proc over a 60 s idle window after a 20 s warm-up. Base binary: 5912 ticks, about 59 CPU-s in 60 s (~98% of a core), which reproduces the fault. Branch binary: 108 ticks, about 1.1 CPU-s in 60 s (~2%). Its first fold took about 3.5 CPU-s and later polls stayed at about 1-2 ticks/s. `rhiza daemon status` showed Last poll advancing, Paused: yes, and the full Skipped list, so the picks and skips were still worked out every poll from the cached fold. The VM load average was 8-20 on 4 cores.

## Links

### Out

- verifies [daemon-idle-poll-no-fold](../capability/rhiza/orchestration/daemon/daemon-idle-poll-no-fold.md)

## History

- `01M6H6SGKR1N5GBG7XPECKN9Z3` Fix The daemon pins a full CPU core: on the 4-core VM rhiza daemon uses ~100% CPU continuously (120 CPU-minutes in 2 hours) because each poll (daemonPoll = 5s, internal/api/daemon.go) takes 10-15 s of work: it re-reads and re-folds the whole log and recomputes picks every time. A poll should be cheap when nothing changed: skip the fold/pick work when the log, git refs, run records and settings are unchanged since the last poll (cache keyed on their stamps), and measure a poll's cost in a test. Agents' rhiza apply calls on the same VM also take minutes under this load
