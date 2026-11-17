# interim-model-deleted-verified

evidence · `01M7K92ENGGRGT632WYVQ9VBDQ` · actor [claude-code-557d88e3](../actor/claude-code-557d88e3.md)

## result

PASS

## statement

At a2ef70341 in a clean detached worktree (load avg 4.5-7.3 on 4 cores), go test -trimpath -count=1 -p 2 -timeout 4m ./internal/api run in two capped halves: -run '^Test[A-Ka-k]' (214 tests) ok in 70.4s and -run '^Test([^A-Ka-k]|$)' (595 tests) ok in 208.0s; go test -list shows the two patterns cover all 809 tests (the other 12 names are benchmarks). Earlier verifiers of this leaf found the acceptance greps, go vet, the projection parity test at 1000 checks, the test counts 2014/398 and make self-check all passing, so the whole acceptance holds.

## Links

### Out

- verifies [interim-model-deleted](../capability/rhiza/views/fast-reads/interim-model-deleted.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
