# writes-from-projection-qa-1

evidence · `01M7KB2196Q8YPX95EFSA2W3KP` · actor [claude-code-49a7d3f2](../actor/claude-code-49a7d3f2.md)

## result

PASS

## statement

QA with the built binary in a scratch repo: add, apply, show, ls, change show, change ls, slice, log, propose, accept and land all read and write through the projection. A dry-run apply whose second write names a parent the first adds, then links it to an existing node, resolved through the pending overlay and wrote nothing (show gizmos not_found afterwards); the real apply then succeeded. Failure: change show of an unknown change set exits 4 'no change set ...' with next rhiza change ls; show nope exits 4; apply with a taken handle and with an unknown parent are refused naming writes[0]. The scoped api tests (TestWorkspace|TestRead|TestChange|TestShip|TestSlice|TestWrite|TestApply|... ) ok in 112 s and the cmd/rhiza scripts apply|change-show|slice|ship-|repair ok. The acceptance greps print nothing.

## Links

### Out

- verifies [writes-from-projection](../capability/rhiza/views/fast-reads/writes-from-projection.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
