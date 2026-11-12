# pr-495-done-needs-a-record

evidence · `01M78779N73Y9M4S5M4NDEK0VQ` · actor [dphaener](../actor/dphaener.md)

## head

e25906c391988be03b48a4a2e3c61892dec35c3f

## pr

https://github.com/dphaener/rhiza/pull/495

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #495 (https://github.com/dphaener/rhiza/pull/495) passed CI at e25906c

## status

verified

## Links

### Out

- verifies [done-refuses-unrecorded-verifier](../capability/rhiza/workflows/workflow-runs/agent-step-done-signal/done-refuses-unrecorded-verifier.md)
- verifies [done-refuses-unrecorded-check](../capability/rhiza/workflows/workflow-runs/agent-step-done-signal/done-refuses-unrecorded-check.md)
- verifies [done-own-command-guide](../capability/rhiza/guide/done-own-command-guide.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
