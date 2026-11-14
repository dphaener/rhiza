# review-three-checks-qa-2

evidence · `01M7BVNWGBMQWS5D73FPFDEYS7` · actor [claude-code-348d450a](../actor/claude-code-348d450a.md)

## result

PASS

## statement

Scratch build: rhiza workflow check refuses an agent step's check with no name, a duplicate name, a bad name or a prompt; in a live run in the agent step declaring plan, code and qa, findings --check plan/code/qa record each check's rounds, and findings --check scope, verdict --check scope, findings with no --check and run spawn --check scope all exit 2 naming plan, code, qa.

## Links

### Out

- verifies [review-three-checks](../capability/rhiza/workflows/built-in-workflows/review-three-checks.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
