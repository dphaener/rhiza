# folds-from-projection-qa-2

evidence · `01M7KGW1EG2Z5AQMJ9FEZ57ACC` · actor [claude-code-69572f11](../actor/claude-code-69572f11.md)

## result

PASS

## statement

Scratch binary built from 96a0eabf3 (go build -trimpath -o $TMPDIR/rhiza-qa2), run in a scratch clone with TMPDIR set to an empty folder, once with no index path (XDG_CACHE_HOME=relative) and once with a read-only cache folder. status --json and inbox --json answer (status counts 14 accepted/40 draft/1 proposed in both modes; inbox items []) and leave 0 rhiza-projection-* folders; during inbox two throwaway folders exist and both are gone at exit. On ~/code/rhiza after one warming run: status 0.06 s, inbox 0.67 s (base build at merge-base 9064447e1: 0.65 s), check 13.45 s and then 9.56 s (base 11.18 s), load about 4 on 4 cores. No regression against base; inbox and check are still over milestone 39's 200 ms/5 s targets, which these leaves do not claim.

## Links

### Out

- verifies [folds-from-projection](../capability/rhiza/views/fast-reads/folds-from-projection.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
