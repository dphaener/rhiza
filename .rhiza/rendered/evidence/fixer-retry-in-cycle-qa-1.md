# fixer-retry-in-cycle-qa-1

evidence · `01M6581N82Q5YSD3TA3CMS96KX` · actor [claude-code-32688ba6](../actor/claude-code-32688ba6.md)

## result

PASS

## statement

At e499bf2 I built rhiza (go build -o $TMPDIR/rhiza-qa ./cmd/rhiza) and ran rhiza run start --workflow review --permission-mode manual in a scratch repo, with a two-check checks workflow, a fake rex that runs the pane command in place and a fake claude: check alpha records one major finding with no leaf in cycle 1, the critic confirms it, and the change-set fixer exits 1 with no commit on its first call and commits with Change:/Review: trailers on its second. The run succeeded (exit 0, review pass, 2 of 2 checks). run status listed two cycle 2 fixer sessions (fail exit 1, then pass with a commit), alpha passed after 2 cycles with no cycle 3, and the fake's log shows fixer call 1 (cycle 2), fixer call 2 (cycle 2), then alpha's check of cycle 2: the cycle's check started only after the retry passed.

## Links

### Out

- verifies [fixer-retry-in-cycle](../capability/fixer-retry-in-cycle.md)

## History

- `01M61GW7TJ6BH6YX1T7SQDWK1J` A failed fixer does not use up a review fix cycle: the fixer is retried within its cycle, and the cycle's checks re-run only once its fixers pass.
