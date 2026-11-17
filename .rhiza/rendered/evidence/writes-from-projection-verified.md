# writes-from-projection-verified

evidence · `01M7K5EP1B9EVRCRZ809CZPQR7` · actor [claude-code-8b64a161](../actor/claude-code-8b64a161.md)

## result

PASS

## statement

At a48d03df4: both acceptance greps (IndexPath/len(w.pending)/newModel(/FoldLog over workspace, read, change, ship, slice, run_model_cache, run_milestone.go; index.Open(ctx.Root, ctx.IndexPath) over read, ship, repair.go) print nothing. go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestWorkspace|...|TestPropertyProjectionParity)' ok (86s); full ./internal/api ok (238s); go vet ./internal/api ok; ./cmd/rhiza -run 'TestScripts/(apply|change-show|slice|ship-|repair)' ok (45s, RHIZA_/SHIP_/REX_ unset). Diff: workspace model/entriesOfSet/changeSets/resolveChangeSet go only through the projection with pending overlaid; openReaderIn, findShipTarget, repair indexBytes take store() from the projection; reserveMilestoneNumber uses w.model; run_model_cache opens via openWriteIndex. TestApplyDryRunResolvesPending covers the pending-overlay apply; TestDefaultViewFoldsNothing checks unknown change set is exit 4 with next rhiza change ls. newModel still defined (model.go:136); no request/response types changed; no TODOs. rhiza ship generate and rhiza render not re-run (they write files; verifier edits nothing).

## Links

### Out

- verifies [writes-from-projection](../capability/rhiza/views/fast-reads/writes-from-projection.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
