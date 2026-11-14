# check-plan-code-topics-qa-2

evidence · `01M7BVNVNXJGVZRGRVERXCHTB8` · actor [claude-code-348d450a](../actor/claude-code-348d450a.md)

## result

PASS

## statement

With a fresh scratch build (go build -trimpath): rhiza guide check-plan, check-code, conduct and review print (3994, 3972, 2690, 1290 bytes); review names check-plan, check-code and check-qa for a conductor run and the older four for ship-change, review and docs-review; check-compliance, check-model and check-quality texts are unchanged from main; rhiza guide install --agent claude into a scratch HOME wrote skills/rhiza-check-plan and rhiza-check-code whose bodies match the topics, and a repeat reported 0 written, 26 unchanged; rhiza guide nosuch exits 4. Acceptance tests pass.

## Links

### Out

- verifies [check-plan-code-topics](../capability/rhiza/guide/agent-guide/check-plan-code-topics.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
