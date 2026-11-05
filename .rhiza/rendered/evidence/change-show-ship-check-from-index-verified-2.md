# change-show-ship-check-from-index-verified-2

evidence · `01M6NZT1GWBB6AEDTJYVPJRVZS` · actor [claude-code-ab9b4f3b](../actor/claude-code-ab9b4f3b.md)

## result

PASS

## statement

Verifier attempt 2 at 07cd7588. go test -trimpath ./internal/api -run TestReadBudget: PASS; change show of an in-flight (draft) and an accepted change set on the 50,000-node model had medians of 30.2/22.5 ms fresh and 17.9/19.5 ms after an unrelated entry, all under 50 ms. go test -trimpath ./internal/api/...: ok. TestChangeShowIndexedMatchesFold (counting log reader, notes included, fresh and after an unrelated entry, and current change set) and TestShipCheckSettledReadsNoLog (ship check --settled, fresh and after an unrelated entry) assert zero log reads: PASS. The four testscripts change-show-serves, ship-check-settles-conflict, ship-check-ci-test-findings and ship-check-follows-moved-head: PASS, unchanged. Re-checked review finding f1 (strict what-if overlay dropped a node moved under one created later): the fix makes newIndexModel use the write path's writeIndex.buildModel view. Ran the finding's repro with a scratch build in a scratch repo (add capability top --record; change new Edge; add early; add late; mv early --parent late; change show --json): leaves [early], matching the fold. The new TestChangeShowIndexedKeepsMovedUnderLater and the 6-seed generated-log TestChangeShowIndexedMatchesFoldGenerated: PASS.

## Links

### Out

- verifies [change-show-ship-check-from-index](../capability/rhiza/views/fast-reads/change-show-ship-check-from-index.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
