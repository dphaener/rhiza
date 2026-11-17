# writes-from-projection-qa-2

evidence · `01M7KGW1Y6MFYC9YV0AZM5WR4F` · actor [claude-code-69572f11](../actor/claude-code-69572f11.md)

## result

PASS

## statement

Scratch binary built from 96a0eabf3 (go build -trimpath -o $TMPDIR/rhiza-qa2), run in a scratch clone with TMPDIR set to an empty folder, once with no index path (XDG_CACHE_HOME=relative) and once with a read-only cache folder. With no index path: apply --dry-run of an add, an add and a link from the second add to the first resolves the link through the pending overlay (rc 0, 2 nodes, 1 link). The real apply writes it and show qa-probe-a --change shows the 'verified by' link. apply of a missing file gives a usage error. show of an unknown handle gives not_found (next rhiza ls). All leave 0 leftovers.

## Links

### Out

- verifies [writes-from-projection](../capability/rhiza/views/fast-reads/writes-from-projection.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
