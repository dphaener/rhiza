# scripts-longest-first-qa-1

evidence · `01M7FPJNSQTG0VBCY68NAWY30H` · actor [claude-code-a30f328a](../actor/claude-code-a30f328a.md)

## result

PASS

## statement

QA at 6e4fed7ee: go test -run '^TestScriptOrder$' passes; with the checked-in timings line 3 made tab-less, real TestScripts fails with 'testdata/script-timings.tsv:3: malformed line' (file restored).

## Links

### Out

- verifies [scripts-longest-first](../capability/rhiza/integrity/self-check/scripts-longest-first.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
