# evidence-commit-survives-moved-head-verified

evidence · `01M7G2CDW5GKSZ0519FDCH910Q` · actor [claude-code-19594e2b](../actor/claude-code-19594e2b.md)

## result

PASS

## statement

At 915b93774, under the machine lock with RHIZA_/SHIP_/REX_ unset: go test -trimpath -p 2 -count=1 ./internal/api -run '^TestCommitOnlyHeadMoved' passed (TestCommitOnlyHeadMoved: hook-injected commit between HEAD read and update-ref; CommitEvidence succeeds, Evidence commit's parent is the racing commit, both on the branch, racer.txt intact, Evidence paths committed and real index clean; TestCommitOnlyHeadMovedEveryTry: fails after 5 tries with update-ref's error). go test -trimpath -p 2 ./internal/api -run 'CommitEvidence|CommitOnly|RunLeaves' passed, -run 'EvidenceCommit|LeafCommit' -count=1 passed, go vet ./internal/api clean, go test -trimpath -p 2 -count=3 ./cmd/rhiza -run 'TestScripts/^headless-run-agent-tabs$' passed all three runs (run once), rhiza ship generate staged 0 files, rhiza render exit 0. The diff touches only internal/api/run_leaves.go, its test and the builder's .rhiza entries; no stubs or TODOs.

## Links

### Out

- verifies [evidence-commit-survives-moved-head](../capability/evidence-commit-survives-moved-head.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
