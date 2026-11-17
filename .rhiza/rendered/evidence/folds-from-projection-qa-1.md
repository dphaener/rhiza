# folds-from-projection-qa-1

evidence · `01M7KB21H3D011SCZ07ZC6A7S9` · actor [claude-code-49a7d3f2](../actor/claude-code-49a7d3f2.md)

## result

PASS

## statement

QA: status, inbox (overlap and blocker items), overlaps, ready, render and check on the scratch repo, and inbox --json identical to the base build and to the no-index run. Failure: after an entry was edited in place, render exits refusing the damaged log with next rhiza repair --dry-run (with and without an index), and repair --dry-run names the entry to restore. Targets on ~/code/rhiza (load ~5-7 on 4 cores, warm, /usr/bin/time, capped scope): status 0.04 s (target 0.2 s, met); inbox 0.53-0.88 s, interleaved with the base build median 0.65 s vs base 0.67 s (no worse than baseline); check 9.7-12.2 s, interleaved 9.7/9.7 s vs base 8.8/11.6 s (over the 5 s target and the 7.1 s baseline under this load, but not worse than the base build). Scoped tests TestStatus|TestInbox|TestUncoveredProposals|TestDaemon*|TestRender ok; scripts status|inbox-|render|daemon-plans-proposals ok.

## Links

### Out

- verifies [folds-from-projection](../capability/rhiza/views/fast-reads/folds-from-projection.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
