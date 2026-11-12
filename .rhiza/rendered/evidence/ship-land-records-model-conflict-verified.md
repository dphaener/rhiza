# ship-land-records-model-conflict-verified

evidence · `01M77PTD0SX7XS0AWZVXTASNAB` · actor [claude-code-89de499f](../actor/claude-code-89de499f.md)

## result

PASS

## statement

Verifier at 5b494d1f6: go test -trimpath ./internal/api/... -run 'Ship|Land' passes (ok), as do ./internal/model/... and TestScripts/describe-golden; gofmt and go vet clean. -v shows TestShipLandRecordsModelConflict (in-process land of a stale-base set exits 3, error carries node/op/other change set/aspects [field:statement], ErrorJSON body carries the same conflict, state records it at head beef000111), TestShipLandHeadBuildRecordsModelConflict (fake head build via ctx.Exec exiting with rhiza land's --json conflict body records the same conflict), TestShipLandClearsModelConflict (missing-dependency refusal clears it), TestShipLandSucceedsAfterRebase (set --rebase, then land succeeds and clears it) and TestShipCheckModelConflict (--model-conflict fails with nothing recorded, passes with intents of both change sets and refused_change_set/other_intent in JSON, passes at 2 fixes, fails at 3 Review:+Model-Conflict: commits, fails once HEAD moved, fails after a missing-dependency refusal), plus TestShipLandWithHeadBuild, all PASS. Diff read: model.StaleBaseError replaces the flat stale-base error (still Unwraps to ErrConflict), checkLand attaches api.ModelConflict, ErrorJSON (used by cmd/rhiza for --json) emits it, headConflict parses the child body's structured conflict rather than headFailure text; a head-build missing-dependency body has no conflict so it clears by the same path (not separately tested). One narrowed rerun once tripped rootguard (a write into ~/code/.rhiza-root, from the live daemon sharing it); three identical reruns passed, and the change touches no root state.

## Links

### Out

- verifies [ship-land-records-model-conflict](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-land-records-model-conflict.md)

## History

- `01M77N1GAGBPBP7EH4EV2CT023` Fix watchdog unmerged PR at step land: run 01M7786QKAZMH7PBJZP7QQA6HX; log: R head's build: conflict: cannot land change set 01M763YX47S5XY0W3P8FVJ54S8: model: conflict: base 01M6RN85C17S5VE1A6SQVD7WE1 of 01M6RN85C17S5VE1A6SNZAXK89 is stale: op 01M77C4T7EBWGBK547B9S31A74 of change set 01M763Z16T0ZGQMY1EC2Y81V0K since wrote field:statement: /home/rhiza/code/rhiza/.git/rhiza/ship/01M763YX47S5XY0W3P8FVJ54S8/rhiza land --change 01M763YX47S5XY0W3P8FVJ54S8 --json: exit status 3
