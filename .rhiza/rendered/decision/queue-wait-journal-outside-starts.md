# queue-wait-journal-outside-starts

decision · `01M79EAN9N4Z20Y6RX8TXBWE8V`

## reasoning

The finding is right: in internal/api/daemon.go (~1020-1037) the ship-owner filter removes candidates before Pick and before the live check, so the daemon never starts an owner-shipped change set and never sees it go live; the Decision's started line would never be written. Limiting the claim to daemon-started runs would leave ship-owner waits, which the intent names as one of the causes (policy), open-ended in the perf analysis, and it would mislabel the accepted-to-ship-start time the perf report (rec 10) measures. Fixing it is cheap and needs no new data: the poll already lists every run journal (rs.List in daemonPoll1), and each journal opens with a TypeStarted entry carrying the change set and the start time (internal/run/engine.go). Taking at from that entry makes the line exact even when the poll runs late or the daemon was down when the run began, and checking it independently of the filters avoids depending on filter order. The rest of queue-wait-journal-file stands: lines only on change, no post-start reasons, no trimming.

## statement

A start the daemon did not make is journaled too. On every poll, before the ship-owner filter and the other pre-start filters, the daemon checks each change set whose last queue-wait line is a pre-start reason: if a run journal under <common>/rhiza/runs/ has a started entry for that change set later than that line, the daemon appends a started line for it. That line's at is the run journal's started time, not the poll's time, and its detail says the run was not started by the daemon (for example 'started outside the daemon'). A started line the daemon writes for its own start keeps detail empty or 'started by the daemon'. So every change set's lines end at started however its run began, including ship-owner change sets and runs started by hand, and the first-line-to-started sum gives the accepted to ship-start time for all of them. The capability queue-wait-journal-lines and its acceptance check gain this case: a change set skipped as ship-owner whose run is then started outside the daemon journals ship-owner then started, with the run's start time.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m79d66a8e9w3e5q9zp0b4wqk-r3-f2](../question/review-01m79d66a8e9w3e5q9zp0b4wqk-r3-f2.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
