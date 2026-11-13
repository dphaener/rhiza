# queue-wait-journal-lines-qa-1

evidence · `01M79WHT4X6BDSMKZ38ZK4DZ98` · actor [claude-code-037e23c0](../actor/claude-code-037e23c0.md)

## result

PASS

## statement

QA cycle 1, built binary, scratch root with three accepted milestone change sets and real rhiza daemon --once polls: two picked ones got started lines and the third got no-slot. A second daemon process wrote no duplicate no-slot. With rhiza.autonomy.ship=owner each change set got a ship-owner line (no-slot to ship-owner and started to ship-owner when no run existed). Under RHIZA_TEST_LOAD1=5, two polls journaled load-high once per change set with the detail '1-minute load 5.00 per core, over the limit of 3.00', then started lines once the load cleared. A change set with a live (hand-written) run journal got no line under ship-owner or load-high. Gap: a run started outside the daemon gets no started line (see the qa finding).

## Links

### Out

- verifies [queue-wait-journal-lines](../capability/rhiza/orchestration/daemon/daemon-picks/queue-wait-journal-lines.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
