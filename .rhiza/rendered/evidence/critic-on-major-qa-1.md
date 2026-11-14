# critic-on-major-qa-1

evidence · `01M7BKZ0AFC3S7025DJDYQ2CQ2` · actor [claude-code-73b338a9](../actor/claude-code-73b338a9.md)

## result

FAIL

## statement

In live runs with a fake conductor, rhiza run review-check from the later gate step failed naming an unjudged blocker and a confirmed major, passed after a downgrade and a clean latest round, and refused a verdict on an unknown finding. But a run whose checks recorded no round passed the gate and succeeded (QA finding f1).

## Links

### Out

- verifies [critic-on-major](../capability/rhiza/workflows/built-in-workflows/critic-on-major.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
