# review-01m79d66a8e9w3e5q9zp0b4wqk-r3-f2

question · `01M79E76C7DF26DQQKYJG6SC34`

## statement

Plan review round 3 left finding f2 (contradiction on queue-wait-journal-file): The Decision says that summing a change set's lines from its first line to started gives the accepted-to-ship-start time, but started is written only when the daemon starts the run. A ship-owner change set is never picked by the daemon (internal/api/daemon.go:1029-1031): the owner starts its run, and every poll after that still skips it as ship-owner, because that filter runs before the live check. Its journal never gets a started line, so its wait looks open-ended, and any other run started by hand has the same gap. The plan should say how a start the daemon did not make gets journaled (for example, a started line when a poll first sees the change set live), or else limit the claim to daemon-started runs. The critic confirmed it: Checked the code: the ship-owner filter (internal/api/daemon.go ~1020-1037) runs before Pick, on every candidate, live or not. An owner-shipped change set is never started by the daemon, so it never gets a started line, and every poll after the owner starts it still skips it as ship-owner. The Decision's claim that the first-line-to-started sum gives accepted-to-ship-start time is false for these change sets and for any run started by hand. The plan has to journal those starts or limit the claim. Major stands. How should it be settled?

## Links

### In

- answers [queue-wait-journal-outside-starts](../decision/queue-wait-journal-outside-starts.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
