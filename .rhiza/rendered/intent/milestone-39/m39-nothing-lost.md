# m39-nothing-lost

intent · `01M7GFT05YPQ76T9QVSC0XS2KG` · parent [milestone-39](../milestone-39.md)

## statement

No test is deleted or skipped to get there: the number of tests go test -list '.*' ./... lists and the number of cmd/rhiza/testdata/script/*.txtar files do not drop below their counts at a98fd4051 (1,948 and 398); a test of deleted interim-model code counts only when a parity case replaces it. Measure: go test -list '.*' ./... | grep -c '^Test' and ls cmd/rhiza/testdata/script/*.txtar | wc -l, on a98fd4051 and on main.

## Links

### Out

- decided-by [projection-targets](../../decision/projection-targets.md)

### In

- serves change set `01M7GFTKCAK0GY3PPDRTPVYF56` (The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).)
- serves change set `01M7GFTS0NVWBXF3CQZD916J1C` (cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
