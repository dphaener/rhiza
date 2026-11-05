# index-catchup-refolds-incrementally-verified

evidence · `01M6NJ5MY4NG83Q6Y8JMEDKNNS` · actor [claude-code-debd4e8b](../actor/claude-code-debd4e8b.md)

## result

PASS

## statement

Verifier claude-code-debd4e8b at 87802e01 ran: go test -trimpath ./internal/index/... -run 'TestPropertyFoldedEveryStep|TestCatchUpFoldsOnlyTouched|TestCatchUpBudget' -v -count=1. Result: PASS. TestPropertyFoldedEveryStep (rapid, 100 cases) appends a generated log's files one at a time in a permuted order. After each catch-up it checks that FoldedOf (history included, reflect.DeepEqual) and the stored fold problems equal changeset.FoldLog over the whole log. The generator supersedes by in-log change sets and by an absent one, so the unknown-supersede clear and re-add paths run. TestCatchUpFoldsOnlyTouched appends one note to set 57 of a 120-set log: read 1, folded 1, decoded 4 (that set's 3 entries plus the new one), no view recomputed, and the result still equals FoldLog. TestCatchUpBudget on the 50,000-node budget model: median 22.6ms after one appended entry (min 15.3ms, max 33.0ms), under the 50ms bar. The full internal/index and internal/changeset packages also pass with -count=1. In the diff, refold folds only the touched sets from decodeEntriesOf, keeps the rest in the folded gob table, merges fold problems per change set and redoes FoldLog's unknown-supersede check in both directions. Deleted files still trigger a full rebuild, and SchemaVersion is bumped to 4.

## Links

### Out

- verifies [index-catchup-refolds-incrementally](../capability/rhiza/views/fast-reads/index-catchup-refolds-incrementally.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
