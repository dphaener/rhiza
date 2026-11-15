# script-timings-file-qa-1

evidence · `01M7FPJNQNXB9EXQYXKFXDMPK8` · actor [claude-code-a30f328a](../actor/claude-code-a30f328a.md)

## result

PASS

## statement

QA at 6e4fed7ee: the checked-in script-timings.tsv has 398 lines, one per *.txtar. make script-timings JSON=<real go test -json of two scripts> rewrote exactly those two lines (restored afterwards); JSON of a build failure, a missing file, neither or both of RUN/JSON exit 2 with the file unchanged. With -file/-scripts on scratch copies: a fake gh printing the replay's script-time lines behind job/step/timestamp prefixes merges, a removed script is dropped, a new one gets 1, a missing previous file is created; a malformed or duplicated previous line, a failing gh and a malformed script-time line exit 1 and leave the file unchanged.

## Links

### Out

- verifies [script-timings-file](../capability/rhiza/integrity/self-check/script-timings-file.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
