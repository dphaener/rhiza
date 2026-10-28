# daemon-resume-once-qa-1

evidence · `01M629Y77SQFP95G5JJS6X96RX` · actor [claude-code-715b2246](../actor/claude-code-715b2246.md)

## result

PASS

## statement

Built rhiza at f087c0e into $TMPDIR/rhiza-qa and drove it by hand in a scratch repo (RHIZA_* env unset, rex shimmed on PATH): a real rhiza run start of ShipB --permission-mode auto (fake claude and rex) was SIGKILLed; rhiza daemon --once then printed 'Resumed run 01M629WTHZX9B3Z7CJ1T7XMNHY' and opened one tab 'run resume <run>'. On a fixture journal with a dead pid, the first --once resumed it and journaled one daemon-resumed entry; the second --once (still interrupted) reported resumed [] and opened no tab; after it finished failed it got no tab.

## Links

### Out

- verifies [daemon-resume-once](../capability/rhiza/orchestration/daemon/daemon-resume-once.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
