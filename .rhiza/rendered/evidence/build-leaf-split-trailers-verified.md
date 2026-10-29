# build-leaf-split-trailers-verified

evidence · `01M62X1TP1FKYBVRD9YNJH4TG5` · actor [claude-code-26f145f7](../actor/claude-code-26f145f7.md)

## result

PASS

## statement

At d773501, ran go test ./internal/api -run 'TestLeafCommit|TestFixLeaf' -count=1 -v: all PASS. Builder half (01c1573): TestLeafCommitEngine runs run.Engine over a real git repo; leaf-split's builder commit (Leaf: in the first trailer paragraph, Co-Authored-By in its own final paragraph) passes with that commit and the verifier starts on it; leaf-split-other (same shape, Leaf: leaf-x) fails with 'new commits name another leaf: leaf-x'. Fixer half (d773501, fixing attempt 1's finding): TestFixLeafEngineSplitTrailers drives the review engine's fixer to commit Change:/Leaf: leaf-a then a separate Co-Authored-By paragraph (asserted in the commit body); the fixer judgement passes on that commit, the verifier starts on it and the step passes. TestFixLeafEngineSplitOtherLeaf: the same split commit naming leaf-x fails the fixer with 'another leaf: leaf-x' and no verifier starts. Read both diffs: test-only changes. make check: fmt, vet and test pass; self-check fails only on stale .rhiza/rendered actor files from this run, unrelated to the leaf.

## Links

### Out

- verifies [build-leaf-split-trailers](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaf-split-trailers.md)

## History

- `01M61GW7BC189FSAN4VEQ1A241` A builder's or fixer's commit counts its Leaf: trailer even when Co-Authored-By sits in a separate final paragraph, and the build and review guides show all trailers in one final paragraph.
