# daemon-picks-qa-1

evidence · `01M629Y68S4X7TF0BQYDNQZYZK` · actor [claude-code-715b2246](../actor/claude-code-715b2246.md)

## result

PASS

## statement

Built rhiza at f087c0e into $TMPDIR/rhiza-qa and drove it by hand in a scratch repo (RHIZA_* env unset, rex shimmed on PATH) over accepted change sets ShipA and ShipB (milestone m, one leaf each), Leafless (milestone, no leaf) and Back (backlog, a leaf): the long-running daemon at default --max-runs started ShipA then ShipB in roadmap order and never Leafless or Back, and its next poll 5s later started nothing more. With journals for ShipB running (a live sleep pid) and ShipA interrupted (dead pid), --max-runs 1 started nothing (B fills the slot, A blocked but takes no slot); with A's run finished failed and B's done, --max-runs 3 started only ShipB, never ShipA.

## Links

### Out

- verifies [daemon-picks](../capability/rhiza/orchestration/daemon/daemon-picks.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
