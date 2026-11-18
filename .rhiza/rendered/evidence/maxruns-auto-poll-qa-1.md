# maxruns-auto-poll-qa-1

evidence · `01M7QZEJ7DYM3VS9V12N1A8703` · actor [claude-code-34e635d3](../actor/claude-code-34e635d3.md)

## result

PASS

## statement

QA by hand with the built binary, three accepted change sets, fake rex and RHIZA_TEST_DETACH_LOG: under auto the first daemon --once started exactly one run (state max_runs 1, max_runs_auto true); with one run journaled the next poll started one more (max_runs 2); with two running, each of load high (RHIZA_TEST_LOAD1/5=3.5), memory low (MEM_AVAILABLE=1000), tmp low (TMP_FREE_BYTES=1000) and disk low (FREE_BYTES=1000) started none with the cap held at 2; maxRuns 2 pinned started none and dropped max_runs_auto; auto with --max-runs 2 pinned 2 (no auto); room again started the third (max_runs 3). TestScripts/daemon-maxruns-auto passes.

## Links

### Out

- verifies [maxruns-auto-poll](../capability/rhiza/orchestration/daemon/daemon-picks/maxruns-auto-poll.md)

## History

- `01M7NVMXZ1YEGZB56ZEAMPX44B` rhiza.daemon.maxRuns=auto lets the daemon start another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit
