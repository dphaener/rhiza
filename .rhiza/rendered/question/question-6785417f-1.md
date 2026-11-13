# question-6785417f-1

question · `01M79D9SYNJBVVE0EG9TRH0RSN`

## session

6785417f-0423-45fb-8791-01a1d2d346d4

## statement

Where does the daemon journal why an accepted change set waits before its run starts? (Today daemon.state.json keeps only the last poll's skip reasons.) / How does the perf analysis read the queue-wait journal? / Is the queue-wait journal trimmed?

## Links

### In

- answers [queue-wait-journal-file](../decision/queue-wait-journal-file.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
