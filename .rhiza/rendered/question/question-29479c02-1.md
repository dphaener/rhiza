# question-29479c02-1

question · `01M6XCDGQSEWXK93F2BV3428TP`

## session

29479c02-8aba-4deb-b5a8-01a1bac2ee1c

## statement

Run 01M6WMX2's queued continuation waited 75 min because the machine cap (maxRuns 3) was held the whole time by three long build/review runs (rhiza X2EMM and X2HMV, nebari WFRNHP). The interrupted run itself took no slot, and no plain start jumped ahead. But the root daemon's round-robin can give a freed slot to another project's plain start ahead of a continuation. Status also shows the queued run only as 'interrupted'. What should the fix do?

## Links

### In

- answers [continuations-first-machine-wide](../decision/continuations-first-machine-wide.md)

## History

- `01M6XC4M3W7Y351WGXE8DTZACT` Fix A queued continuation can wait far longer than its slot should take: run 01M6WMX236BTT04TYQDACGKBZ4 (autonomy per workflow, PR #376 green) was queued with rhiza run resume at 2026-11-08T08:09Z and was still not started at 09:24Z; daemon status listed it as 'interrupted' among the live runs next to 3 running ones at rhiza.daemon.maxRuns 3, while two plain starts were 'waiting for a slot'. Check whether an interrupted run with a queued continuation counts against the cap itself (so it needs live <= cap-1 to start), and whether continuations take the next free slot ahead of new starts as #226/#243 intend. Unstuck by driving it with rhiza run resume --foreground.
