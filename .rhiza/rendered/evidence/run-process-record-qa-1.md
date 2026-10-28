# run-process-record-qa-1

evidence · `01M629Y87P3FBS3ER8S4YDFJ6M` · actor [claude-code-715b2246](../actor/claude-code-715b2246.md)

## result

PASS

## statement

Built rhiza at f087c0e into $TMPDIR/rhiza-qa and drove it by hand in a scratch repo (RHIZA_* env unset, rex shimmed on PATH): rhiza run start of ShipB --permission-mode auto with REX_SESSION/REX_BLOCK=block:fakeb (fake claude and rex) journaled {type:process, host Mac.localdomain, pid 11821 (the run's own pid), rex_block block:fakeb} and called 'rex block inspect -s session:fake --json block:fakeb' (the shim's reply did not match Rex's shape, so no window was recorded). After SIGKILL, rhiza daemon status listed that run as interrupted and a fixture run with a live sleep pid as running.

## Links

### Out

- verifies [run-process-record](../capability/rhiza/workflows/workflow-runs/run-process-record.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
