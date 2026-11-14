# skill-check-verifier-check-qa-1

evidence · `01M7E24TPZ417QRQ1XYMVWE1EB` · actor [claude-code-98071bd0](../actor/claude-code-98071bd0.md)

## result

PASS

## statement

QA: scratch build with fakes. --role verifier: good passes 3 of 3, exit 0, no builder session launched; one writing no Evidence fails evidence (and session-done) by name, exit 1; crashing verifier fails, exit 1. --role check: good records a round and passes 3 of 3; bare and crashing checks fail findings by name, exit 1 (re-run for 3 cycles, as designed). --keep names a scratch repo outside the current repo; the user repo is unchanged after every case. TestScripts skill-check-roles passes.

## Links

### Out

- verifies [skill-check-verifier-check](../capability/rhiza/workflows/workflow-authoring/skill-check-verifier-check.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
