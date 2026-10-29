# commit-trailers-every-paragraph-verified

evidence · `01M62WEJ6JFF50BR7KH9EH2FMC` · actor [claude-code-c3b383c8](../actor/claude-code-c3b383c8.md)

## result

PASS

## statement

At b39b297 I ran go test ./internal/api -run 'TestParseTrailerLog|TestGitCommits|TestLeafCommit' -count=1 -v and all passed. TestParseTrailerLog (unit test of parseTrailerLog, the reader used by check's gitCommits and by the runner's LeafCommit/ReviewCommit) finds Change: CS4, Actor: a, Leaf: split and Review: R1 when Co-Authored-By sits in a separate final paragraph (case ddd), and still reads Change:/Leaf: when all trailers share the final paragraph (case eee). The subject line is not read as a trailer (fff). TestLeafCommit/'Co-Authored-By in a later paragraph' finds the Leaf: trailer on a real git commit built with split paragraphs. In a scratch worktree I restored the old %(trailers) check.go from edb268a and the same tests failed (TestParseTrailerLog, TestLeafCommit), so the tests pin the fix. I read the diff: trailerLogFormat now emits %B and parseTrailerLog scans every body line after the subject. make check passed fmt-check, vet and test; self-check failed only on stale .rhiza/rendered actor/index files from this run's own model writes, not on this leaf's code.

## Links

### Out

- verifies [commit-trailers-every-paragraph](../capability/rhiza/integrity/commit-trailers-every-paragraph.md)

## History

- `01M61GW7BC189FSAN4VEQ1A241` A builder's or fixer's commit counts its Leaf: trailer even when Co-Authored-By sits in a separate final paragraph, and the build and review guides show all trailers in one final paragraph.
