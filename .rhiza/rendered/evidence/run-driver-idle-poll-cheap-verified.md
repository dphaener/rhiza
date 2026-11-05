# run-driver-idle-poll-cheap-verified

evidence · `01M6NX7PPP8AFYX0HXPB78XSMQ` · actor [claude-code-d39564a9](../actor/claude-code-d39564a9.md)

## result

PASS

## statement

Verifier at de45583d. Read the diff: stepRunner's Leaves, LeafAnswer, Answer and Evidence (the calls the leaves and checks steps and their arbiter waits make) now go through readModel (internal/api/run_model_cache.go), which keeps the worktree's index open, refreshes it by a stat fingerprint (Index.Refresh/Stamp) and keeps answers while the stamp, or the per-change-set Signature, is unchanged. Ran go test -trimpath -count=1 -v -run TestDriverPollIsCheap ./internal/api/ (PASS; also PASS under -race): on a gen'd log of 120 landed + 6 in-flight change sets with IndexPath set and a live session record, 5 idle polls show ctx.readLog 0, catch-up Read/Decoded/Folded 0 and no answer recomputed; after an entry in an unrelated change set the next poll reads the log 0 times and folds exactly 1 change set; after an entry in the run's change set the leaves go 2->3, and a session done mark is seen (SaidDone). Ran go test -trimpath -count=1 -v -run TestReadBudget ./internal/api/ (PASS, quiettest, medians): driver poll idle median 4.06ms CPU (budget 5ms), after an unrelated entry 5.51ms (budget 50ms); a second run (whole package, go test -trimpath -count=1 ./internal/api/, ok) gave 2.24ms and 6.33ms. go test -trimpath ./internal/index/... ok, incl. TestSignature.

## Links

### Out

- verifies [run-driver-idle-poll-cheap](../capability/rhiza/orchestration/daemon/run-driver-idle-poll-cheap.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
