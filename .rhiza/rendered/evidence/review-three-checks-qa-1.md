# review-three-checks-qa-1

evidence · `01M7BKYZYWE18GAFN1N7M6RCVM` · actor [claude-code-73b338a9](../actor/claude-code-73b338a9.md)

## result

PASS

## statement

With the built binary in a scratch repo: rhiza workflow check refuses an agent step's check with no name, a duplicate name, a bad name, or a prompt, and still refuses check on a command step and critic on an agent step; in a live run with a fake claude in an agent step declaring plan, code and qa, rhiza run findings --check plan/code/qa recorded each check's round, and --check scope exited 2 naming plan, code, qa.

## Links

### Out

- verifies [review-three-checks](../capability/rhiza/workflows/built-in-workflows/review-three-checks.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
