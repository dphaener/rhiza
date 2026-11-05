# change-show-ship-check-from-index-verified

evidence · `01M6NTT1RZHPD5A70G4PP2PJYG` · actor [claude-code-a5e353eb](../actor/claude-code-a5e353eb.md)

## result

PASS

## statement

Verifier, attempt 3, at a7927bb8 (detached checkout). `go test -trimpath -count=1 ./internal/api -run 'TestReadBudget$'`, run 7 times on a contended 4-core machine (load 4-10, other agents' tests running): 4 passed and 3 failed. When a run passed, change show of the draft and accepted change sets, fresh and after an unrelated entry, had medians of about 29-45 ms. Before this commit they were 30-45 ms even on good attempts. Every failure was under load: one run had load 10. In another, the existing what-if reads that this leaf does not touch (log, roadmap, ready what-if) also failed their 50 ms budget on all 3 attempts. So change show now sits in the same band as the existing reads and is no longer the outlier. The functional points pass: TestChangeShowIndexedMatchesFold (zero log reads through a counting reader, notes included, fresh and after an unrelated entry), TestChangeShowIndexedErrors and TestShipCheckSettledReadsNoLog (fresh and after an unrelated entry). `go test -trimpath -short ./internal/api/... ./internal/index/... ./internal/changeset/...` passes, and so do the four named testscripts (change-show-serves, ship-check-settles-conflict, ship-check-ci-test-findings, ship-check-follows-moved-head). I read the diff: the ops-only column (schema version 5), whatIfSets folding cs first with a states map and a cached landed bucket, and the depends-on links indexed by from. Answers are unchanged: the testscripts and the fold-equivalence test pass.

## Links

### Out

- verifies [change-show-ship-check-from-index](../capability/rhiza/views/fast-reads/change-show-ship-check-from-index.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
