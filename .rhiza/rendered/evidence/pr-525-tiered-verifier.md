# pr-525-tiered-verifier

evidence · `01M79MAGE5YJ8MMWJN37QJW3ZP` · actor [dphaener](../actor/dphaener.md)

## head

8c1ac365f0bbde5cda7494331e9b6594f9fa0b83

## pr

https://github.com/dphaener/rhiza/pull/525

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #525 (https://github.com/dphaener/rhiza/pull/525) passed CI at 8c1ac36

## status

verified

## Links

### Out

- verifies [verifier-tier-fields](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-fields.md)
- verifies [verifier-tier-choice](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-choice.md)
- verifies [verifier-tier-defaults](../capability/rhiza/workflows/built-in-workflows/build-workflow/verifier-tier-defaults.md)
- verifies [build-workflow-file](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-workflow-file.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
