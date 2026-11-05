# run-driver-idle-poll-cheap-verified-2

evidence · `01M6P03YK2CKY5Q7HGACHP9CVE` · actor [claude-code-00b28496](../actor/claude-code-00b28496.md)

## result

PASS

## statement

At a9374a62 (attempt 2, fixing review finding f2): read the diff. Evidence is now cached under the run's change set (not the zero id) and over the index judges the creator by FoldedOf of only the change set holding the Evidence node's create (writeModel.creator); the log path still folds every set lazily. go vet ./internal/api passes; go test -trimpath -count=1 ./internal/api/... passes (ok, 427s). Verbose run: TestDriverPollIsCheap PASS on a generated log of 120 landed + 6 in-flight change sets with the index: idle polls catch up nothing, compute no answer and read the log 0 times; after an unrelated change set's entry the poll folds 1 change set, recomputes 0 answers (Evidence included) and reads the log 0 times; a write to the run's change set and a session done mark are both seen. TestReadBudget (quiettest, 50,000-node model) driver poll, now judging a verifier's Evidence each poll: idle median 3.086ms CPU (budget 5ms), after an unrelated entry median 3.267ms CPU (budget 50ms). TestEvidenceCommit/TestEvidenceCommitNone PASS.

## Links

### Out

- verifies [run-driver-idle-poll-cheap](../capability/rhiza/orchestration/daemon/run-driver-idle-poll-cheap.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
