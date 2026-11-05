# pr-271-fix-test-temp-folders-leak-and

evidence · `01M6MMTMEAZRBXHP4AP08GT40G` · actor [dphaener](../actor/dphaener.md)

## head

4e7621e5b49b8c027b73e77d34aad2daa57606c5

## pr

https://github.com/dphaener/rhiza/pull/271

## result

1 check passed: check SUCCESS

## statement

Pull request #271 (https://github.com/dphaener/rhiza/pull/271) passed CI at 4e7621e

## status

verified

## Links

### Out

- verifies [test-temp-folders-swept](../capability/rhiza/integrity/self-check/test-temp-folders-swept.md)

## History

- `01M6M1M10N7AP6F0DHGANKXEA1` Fix Test temp folders leak and fill the disk: TestMain in internal/api/budget_test.go removes the generated rhiza-api-budget-* model (~146MB) only when the test binary exits normally, so every go test timeout (common under load) leaves one behind; quiettest leaves rhiza-quiet-* the same way, and slice/index budgets their own. On the VM 124 budget folders (~18GB) and 369 quiet folders piled up in TMPDIR by 2026-11-04 and the disk fell to 19GB free. Have each package remove its own stale folders (older than an hour, by prefix) when it starts, or put them under one per-run directory make test deletes, and test it
