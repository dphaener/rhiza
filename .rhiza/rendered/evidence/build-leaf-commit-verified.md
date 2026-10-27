# build-leaf-commit-verified

evidence · `01M5ZMT5GWE84SZR5EPY19KX6V` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

TestLeafCommitEngine runs the engine against a real (hermetic) git repository over three leaves: leaf-none's builder makes no commit and its attempt fails with 'no new commit since <rev>'; leaf-other's builder commits naming leaf-x and fails with 'new commits name another leaf: leaf-x'; neither starts a verifier; leaf-own's builder commits with its Leaf: trailer, passes with that commit, and its verifier's leaf-session and LeafSpec carry that commit and RHIZA_COMMIT. TestLeafCommit covers the trailer search cases; TestLeavesStepLeafCommit the engine side with a fake. All passed.

## tests

go test -count=1 ./internal/api -run 'TestLeafCommit$|TestLeafCommitEngine'; ./internal/run -run TestLeavesStepLeafCommit (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-leaf-commit](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaf-commit.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
