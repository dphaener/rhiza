# queue-wait-journal-file

decision · `01M79DCE6ZA8ZR6ANCNP74XKJ8`

## reasoning

Today daemon.state.json is rewritten whole every poll, and its Skips hold only the last poll's reasons (internal/daemon/state.go), so a wait's history needs its own file. The repo already keeps an append-only JSONL journal in the daemon folder: usage.limits.jsonl, written with one O_APPEND write per line (internal/daemon/usage.go). The run journals are JSONL under <common>/rhiza/runs/ too. Putting queue.wait.jsonl next to them follows that precedent and keeps it per project and out of .rhiza/, so it adds no model file format. The skip reasons are already fixed constants in internal/daemon/pick.go, and the poll already works them out in daemonPoll1, so journaling them means writing down what the poll knows. The machine (load) hold needs no new reason: under a machine hold daemonPoll1 already records a Skip for each candidate with the hold's reason, load-high or memory-low (internal/api/daemon.go, internal/daemon/load.go), which covers the intent's 'load hold'. The 2026-11-12 perf report (~/orchestrator/perf-2026-11-12) was an ad hoc analysis of journal files on disk, not a rhiza command. Reading this file the same way meets 'the perf analysis can read it' with the smallest change, and a typed reader in internal/daemon gives the acceptance check something to test. Writing only on change keeps the file small, so trimming is not needed. Not trimming also keeps the ALL-history comparison the report depends on (60m ALL against 27m RECENT).

## statement

The daemon appends why an accepted change set waits to an append-only JSON-lines file, queue.wait.jsonl, in the daemon's folder under the project's git common dir (<common>/rhiza/, beside daemon.state.json and usage.limits.jsonl). It writes a line only when a change set's wait changes, not every poll: when a poll first skips it, when its reason changes, and when its run starts. Each line has at, change, title, reason and detail, where reason is the poll's pre-start Skip reason (no-slot, not-ready, no-milestone, no-leaves, disk-low, usage-reserve, paused, ship-owner, workflow-unresolved, and load-high or memory-low for the machine hold), or started when the daemon starts the run. The post-start Skip reasons (live, held, unjournaled, failed, shipped) belong to a change set whose run has already started and are never journaled, so a change set's lines end at started until it waits to start again. The fields are a fixed type in internal/daemon with a reader (one function that returns every line, oldest first, and skips lines that do not parse), as LimitHit and ReadLimitHitsIn are. The perf analysis reads this file directly, as it reads the run journals under <common>/rhiza/runs/. Summed per change set, from its first line to started, it gives the accepted to ship-start time split by reason. This change set adds no CLI command or --json shape for it. The journal is not trimmed: because it only records changes, it holds a few short lines per change set, and the perf report compares RECENT with ALL history, which trimming would lose.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-6785417f-1](../question/question-6785417f-1.md)
- derived-from [queue-wait-journal-plan](../source/queue-wait-journal-plan.md)

### In

- decided-by [queue-wait-journal-reader](../capability/rhiza/orchestration/daemon/daemon-picks/queue-wait-journal-reader.md)
- decided-by [queue-wait-journal-lines](../capability/rhiza/orchestration/daemon/daemon-picks/queue-wait-journal-lines.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
