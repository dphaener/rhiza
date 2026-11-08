# queued-continuation-slot-interview

source · `01M6XCJ8RAFR0TEKRWQF04460T`

## quotes

```json
[]
```

## raw_sha256

65e3a79769e8ed23612ddeecde8691f8b081f2bc69b89e0b01daf1bf21f10a3a

## session

29479c02-8aba-4deb-b5a8-01a1bac2ee1c

## summary

Autonomous plan run, no owner. The planner looked into why run 01M6WMX236BTT04TYQDACGKBZ4's queued continuation (rhiza run resume, 2026-11-08 08:09Z) did not start until it was driven by hand at 09:24Z. Run journals show the machine cap (maxRuns 3) was held the whole time by three long runs: rhiza 01M6X2EMMHKMQ9W4TDJH3ERXSD and 01M6X2HMV36WFWZYPMWK8XKFXG, and nebari 01M6WFRNHPRBVRDQ4KJHB6CSR2. No run started or ended in that window. The code agrees: an interrupted run takes no slot (daemonDemand and daemon.Pick count only LiveRunning), and within a project continuations are kept a slot ahead of ship picks. Two gaps remain. daemon.RoundRobin takes each project's ready list in rotation, so a freed machine slot can go to another project's plain start or pick ahead of a queued continuation. And daemon status lists the queued run only as 'interrupted'. The planner raised one question on scope (priority plus visibility, one over the cap, or no cap). The arbiter ruled in Decision continuations-first-machine-wide: continuations first machine-wide, then queued starts and plan requests, then ready change sets, within maxRuns and per-project caps; daemon status shows a queued continuation as waiting for a slot with its place in line. No change to counting, no pre-emption, no slot outside maxRuns.

## title

Queued continuation slot wait, autonomous plan interview

## Links

### In

- derived-from [continuations-first-machine-wide](../decision/continuations-first-machine-wide.md)

## History

- `01M6XC4M3W7Y351WGXE8DTZACT` Fix A queued continuation can wait far longer than its slot should take: run 01M6WMX236BTT04TYQDACGKBZ4 (autonomy per workflow, PR #376 green) was queued with rhiza run resume at 2026-11-08T08:09Z and was still not started at 09:24Z; daemon status listed it as 'interrupted' among the live runs next to 3 running ones at rhiza.daemon.maxRuns 3, while two plain starts were 'waiting for a slot'. Check whether an interrupted run with a queued continuation counts against the cap itself (so it needs live <= cap-1 to start), and whether continuations take the next free slot ahead of new starts as #226/#243 intend. Unstuck by driving it with rhiza run resume --foreground.
