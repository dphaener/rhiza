# queue-wait-journal-lines-qa-2

evidence · `01M79XSHA210BK2MP3NPCTBCEA` · actor [claude-code-72ff73af](../actor/claude-code-72ff73af.md)

## result

PASS

## statement

QA cycle 2 at 6f0d5d935, built binary (go build -trimpath ./cmd/rhiza) in a hermetic scratch Rhiza root (env -i, fake rex). With rhiza.autonomy.ship=owner one rhiza daemon --once journaled ship-owner for three accepted milestone change sets a, b, c. Started a's run by hand (rhiza run start --foreground, fake rex, stopped by timeout): the next two polls wrote exactly one line for a, reason started, at 07:39:55.043401079Z equal to the run journal's started entry, detail 'started outside the daemon'; b and c got nothing new. Under a real load-high hold b and c moved to load-high with the load detail and a (interrupted run) got no line. With the load guard off and one slot, the daemon resumed a's run and b and c went to no-slot. With three slots it started b and c and wrote started lines with an empty detail. Two more polls, after their run journals appeared (started a fraction of a second after those lines), wrote nothing (12 lines in all). go test -trimpath -count=1 ./internal/api/... -run QueueWait: all four QueueWait tests PASS, including TestQueueWaitShipOwnerStartedOutside.

## Links

### Out

- verifies [queue-wait-journal-lines](../capability/rhiza/orchestration/daemon/daemon-picks/queue-wait-journal-lines.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
