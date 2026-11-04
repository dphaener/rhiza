# roadmap-paged-verified

evidence · `01M6J0A23JR389EM5T2489WTAF` · actor [claude-code-e9b145b3](../actor/claude-code-e9b145b3.md)

## result

PASS

## statement

Verifier claude-code-e9b145b3 at b0414782 (diff ef89d415..b0414782 read). go test ./cmd/rhiza -run TestScripts/(roadmap-paged|describe-golden) passes: roadmap-paged.txtar builds m1{A,B}, m2{A}, backlog{C}, walks --limit 1 pages and cmps each page's group and change set against m1 A / m1 B / m2 A / Backlog C (so A appears once under each milestone, every entry once) and the kinds milestone x3 then backlog; the second page repeats the split m1 heading; the walked entries cmp equal to --all; roadmap --limit 1 --json and api call roadmap {"limit":1} cmp equal; --limit -1, --limit 1001 and --cursor !! each exit 2. describe-golden passes with limit/cursor/all and next_cursor. TestRoadmapOverAPageIsReadWhole (api: 104 change sets; default page short with cursor, --all whole, healthRoadmap counts all), roadmapapp TestReadsARoadmapOverAPageWhole and modelapp TestRoadmapTabReadsARoadmapOverAPageWhole all pass. All three attempt-1 findings are covered. fmt-check ok; rhiza check on a clean checkout of the commit ok; go test ./cmd/rhiza with RHIZA_* unset and -timeout 30m ok. Only wall-clock budget tests failed (modelapp TestGeneratedModel, api TestReadBudget, slice TestBudget) under load average 24-29 on 4 cores; TestGeneratedModel misses its 50ms budget the same way at the pre-leaf base (112-147ms vs 85-162ms at HEAD), so the failures come from load, not this leaf.

## Links

### Out

- verifies [roadmap-paged](../capability/rhiza/views/paged-lists/roadmap-paged.md)

## History

- `01M6H3A65RDP761HKE82STZSQZ` log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan)
