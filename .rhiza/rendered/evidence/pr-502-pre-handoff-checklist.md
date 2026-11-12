# pr-502-pre-handoff-checklist

evidence · `01M78M4S0RY4KYT7RGHRT19RM1` · actor [dphaener](../actor/dphaener.md)

## head

4fbae53b1d87dd8ab07b06f27048f601ac523602

## pr

https://github.com/dphaener/rhiza/pull/502

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #502 (https://github.com/dphaener/rhiza/pull/502) passed CI at 4fbae53

## status

verified

## Links

### Out

- verifies [builder-prehandoff-checklist](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/builder-prehandoff-checklist.md)
- verifies [fixer-prehandoff-checklist](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/fixer-prehandoff-checklist.md)
- verifies [verifier-reruns-reported-checks](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/verifier-reruns-reported-checks.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
