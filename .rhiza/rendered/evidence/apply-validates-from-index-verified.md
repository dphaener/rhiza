# apply-validates-from-index-verified

evidence · `01M6NSA5M87SMXSKH7PVQBE9QM` · actor [claude-code-766eb5c6](../actor/claude-code-766eb5c6.md)

## result

PASS

## statement

Verified at 7457eacf (attempt 2). Read the diff: the write view (internal/api/write_indexed.go writeIndex.model) now loads only the as-is rows the in-flight ops touch, applies the dependencies' ops in dependencyOrder and then the change set's own as newModelSets/model.apply do, and answers through index.OverlayOf. No log fold. (1) go test -trimpath -count=1 ./internal/api -run 'TestWritesThroughIndexMatchFold|TestApplyIndexedReadsNoLog' -v: both PASS. TestApplyIndexedReadsNoLog uses a counting log reader and a fresh index and asserts zero log reads for an apply of several writes. TestWritesThroughIndexMatchFold now generates the log (genLog: 3-5 random change sets with random writes, random depends-on links, end states draft/proposed/accepted/landed/withdrawn) and then random add/set/mv/retire/link/unlink/apply writes to random change sets, with no generator exclusions. Every write got the same error or the same ops through the index as through the fold, with 18-29 accepted and 51-62 refused per seed. This answers finding f2. In a scratch copy I raised it to 16 seeds and all passed (13-29 accepted per seed). (2) Finding f1 re-probed in a scratch copy with an extra test: indexedLog draft, retire ssd, then dry-run set/mv/link/retire of nvme (ssd's child). The indexed and folded answers are byte-identical, all accepted, and the indexed path read the log 0 times. f1 is fixed. (3) go test -trimpath -count=1 ./internal/api -run 'TestReadBudgetApply$' -v, twice, under load 5-7 on 4 cores: one apply write on the 50,000-node model had medians of 131.5ms and 102.5ms (budget 200ms) and read the log zero times; PASS. `-run TestReadBudget` also matches TestReadBudget itself. Its change show/ready timings went over 50ms under this load. Those belong to the change-show-ship-check-from-index leaf, already in its findings, and are not this leaf's point. (4) go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(apply|link)$': PASS. (5) go test -trimpath -count=1 -short ./internal/api/... ./internal/index/...: ok.

## Links

### Out

- verifies [apply-validates-from-index](../capability/rhiza/views/fast-reads/apply-validates-from-index.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
