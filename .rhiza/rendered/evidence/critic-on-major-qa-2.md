# critic-on-major-qa-2

evidence · `01M7BVNW822YSGWZP2NRQ43Y20` · actor [claude-code-348d450a](../actor/claude-code-348d450a.md)

## result

PASS

## statement

Live runs (conduct agent step declaring plan, code, qa, then a gate command step running rhiza run review-check): all three clean rounds pass and the run succeeds; no rounds at all, or qa never spawned, fails naming the missing check(s) (cycle-1 f1 fixed); a code blocker with no critic fails naming 'code: f1 blocker'; a critic downgrade to minor passes; a confirmed blocker fails; a major in round 1 then a clean round 2 passes and the reverse fails (latest round only); review-check --findings with no findings exits 3.

## Links

### Out

- verifies [critic-on-major](../capability/rhiza/workflows/built-in-workflows/critic-on-major.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
