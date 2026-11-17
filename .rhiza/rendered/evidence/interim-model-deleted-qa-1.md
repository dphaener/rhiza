# interim-model-deleted-qa-1

evidence · `01M7KB21MKW682QWWPAY1V6EV0` · actor [claude-code-49a7d3f2](../actor/claude-code-49a7d3f2.md)

## result

PASS

## statement

QA: every acceptance grep prints nothing except write_indexed.go:61 (the throwaway-index check) and render.go's one FoldLog (the HoldsLog fallback). go test ./internal/api ran past my 240 s cap under load ~7 with no output, so I ran it in two halves: -run '^Test[A-Q]' ok in 110 s and -run '^Test[R-Z]' ok in 165 s (no failure, no hang). The parity test with -rapid.checks=1000 ok in 63 s. The behaviour checks above (the no-index reads, the damaged index rebuilt, inbox and status matching the base build byte for byte) show no fold fallback is left and no --json shape changed.

## Links

### Out

- verifies [interim-model-deleted](../capability/rhiza/views/fast-reads/interim-model-deleted.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
