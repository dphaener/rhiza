# byo-failed-run-retry-qa-1

evidence · `01M7REYV8SDTPW4W85C0JYRKKG` · actor [claude-code-8a84713f](../actor/claude-code-8a84713f.md)

## result

PASS

## statement

QA at e288b6b36 with a scratch build: close-out failing once; daemon --once queued a --retry continuation (failed retry 1 of 1), the next poll resumed it, the resume retried only close-out (visit 2, retried) and the run succeeded with intervention_count 0. With close-out always failing, the retry failed and the next poll left it to the inbox with no second retry. TestScripts/byo-workflow-retry passes.

## Links

### Out

- verifies [byo-failed-run-retry](../capability/rhiza/workflows/workflow-runs/byo-failed-run-retry.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
