# ci-speed-shape

decision · `01M6Q5FFGR7Q2E5A4YRAJF88P1`

## reasoning

Read the question, the change set's intent (owner: 'make it so', shard TestScripts e.g. 4 by script name, budgets as their own job, cache modules, keep make check full locally, one aggregate status so pr-green/ci-fix keep working), .github/workflows/ci.yml (one job 'check' running make check), the Makefile's test target (rest of packages in background, cmd/rhiza under the lock, budgets with -p 1 after), cmd/rhiza/main_test.go (testscript.Run over Dir; go-internal v1.14.1 supports Params.Files), and internal/api/ship.go (parsePRView reads every statusCheckRollup entry; checksVerdict needs all completed and SUCCESS/NEUTRAL/SKIPPED, so a skipped aggregate would pass: hence if: always()). Intent: the roadmap's open milestones are 27 (Nebari import) and 28 (web UI), neither about CI; milestone-25 is the standing home for owner-requested fixes and backlog work and recent fix change sets (01M6PQTK) serve it, so it fits best without inventing a milestone (planning milestones needs the owner). Shards: round-robin over sorted names is deterministic, needs no timing data to maintain, picks up new scripts automatically and keeps local runs unchanged when the variable is unset; testscript's Files param makes it a small change. internal/api: at 280 s it alone breaks the ~5 min goal once setup is added, so it must be split; test-name sharding via go test -list applies the same round-robin idea without touching its 90 test files. Verdict: one workflow keeps one run ID for log-based ci-fix, fail-fast off surfaces all failures in one cycle, and an always-run aggregate named like today's job keeps a single required status.

## statement

(1) Intent: change set 01M6Q58 serves milestone-25 (01M6HKVRE1, 'The backlog is empty: every change set planned from the owner's queued requests and the backlog is built and shipped by the daemon on its own'), as the other owner-requested fixes this week do (e.g. 01M6PQTK, the owner-run plan gate fix, serves it); link 01M6Q58 serves milestone-25 and no other milestone. (2) Script shards: TestScripts in cmd/rhiza/main_test.go reads RHIZA_TEST_SHARD=i/n (1-based, e.g. 2/4); when set it lists testdata/script/*.txtar, sorts the names, keeps every name whose index mod n is i-1 (round-robin, so the slow families such as daemon-* and run-* spread across shards) and passes them to testscript.Params.Files instead of Dir; unset, it runs every script as today, so make check and make test locally are unchanged. A script added or renamed lands in a shard automatically; nothing lists scripts by hand. CI runs a matrix of 4 jobs, each `go test -trimpath -skip '<BUDGET_TESTS>' ./cmd/rhiza` with its RHIZA_TEST_SHARD. (3) internal/api: it gets the same treatment as a second matrix dimension rather than staying in the 'rest' job: CI splits internal/api's top-level tests (from `go test -list . ./internal/api`, budget tests excluded) round-robin into 2 shards passed as an anchored -run regex, each its own job; every other package runs in one 'rest' job (minus the budget tests), the timing budgets (BUDGET_TESTS over BUDGET_PKGS, -count=1 -p 1) run in their own job, and fmt-check, vet and self-check run in a 'lint' job. Put the shard selection in Makefile targets (e.g. make test-scripts SHARD=i/4, make test-api SHARD=i/2, make test-rest, make test-budgets) so CI calls make, not ad-hoc go commands, and make test/make check keep running everything. If internal/api shards still exceed about 3 min of test time after the build cache, raise its shard count; do not drop or weaken its tests. (4) One verdict: keep all jobs in the single CI workflow (one Actions run, so ship check --test-failures and ci-fix keep one run ID to read logs from), set strategy.fail-fast: false so every failing shard reports its failures for ci-fix, and add an aggregate job named `check` (the existing job's name, so any required status keeps matching) with needs: on every other job and if: always(), which fails unless every needed job's result is success. It must run under if: always() because ship check's checksVerdict counts SKIPPED as passing, so a skipped aggregate would read green. checksVerdict already requires every rollup check completed and passing, so the shard jobs and the aggregate together give pr-green one settled verdict; ship check needs no change unless a test shows otherwise, and a testscript or unit test over parsePRView/checksVerdict covers a rollup of several jobs with one failed shard and a cancelled aggregate. Cache Go modules and the build cache with actions/setup-go's cache keyed on go.sum in every job.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-c92d1843-1](../question/question-c92d1843-1.md)
- supersedes [t1-d12-ci-one-job](t1-d12-ci-one-job.md)
- derived-from [ci-speed-plan-interview](../source/ci-speed-plan-interview.md)

### In

- decided-by [test-scripts-sharded](../capability/rhiza/integrity/self-check/test-scripts-sharded.md)
- decided-by [make-test-shard-targets](../capability/rhiza/integrity/self-check/make-test-shard-targets.md)
- decided-by [ci-yml-parallel-jobs](../capability/rhiza/integrity/self-check/ci-yml-parallel-jobs.md)
- decided-by [ship-check-multi-job-verdict](../capability/rhiza/integrity/self-check/ship-check-multi-job-verdict.md)

## History

- `01M6Q58T8EG0QY7YDQNWHW6NEA` CI takes 11-18 minutes per PR; get it to about 5 or less. Measured on run 42136754545 (2026-11-06, 17.4 min total): 'Run checks' is 1046 s, of which the cmd/rhiza testscripts take 655 s (one package, scripts run serially under the test lock) and are the whole critical path; internal/api takes 280 s in parallel; the timing-budget phase after the suite takes about 3 min (-p 1), including 152 s of internal/api; and 34 'go: downloading' lines show no Go module cache. Shard TestScripts across parallel CI jobs (e.g. a matrix of 4 by script name), run the timing budgets as their own parallel job (only the budget tests, not the whole package), and cache Go modules and the build cache (actions/setup-go cache or actions/cache keyed on go.sum). Keep make check meaning the full suite locally and keep one required status for the PR (an aggregate job), so ship-change's pr-green and ci-fix logic, which read the PR's checks, keep working; update ship check if needed. Owner, 2026-11-06: "11-18 minutes for CI? That's absurd. We need to fix that." and "make it so".
