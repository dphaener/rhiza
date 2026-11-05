# test-temp-cleanup-plan

source · `01M6M1RQZAQ2KV2FE01W9RXWHH`

## quotes

```json
[]
```

## raw_sha256

ca3855178467b484cdd138ccf26af272efa7c867142860d2b3fb1efcc9239c02

## session

45b188d1-8ca7-4482-a8d3-01a1a81a2270

## summary

Autonomous plan run for the change set fixing leaked test temp folders (rhiza-api/index/slice-budget-*, rhiza-quiet-*) that filled the VM disk when go test timeouts skipped TestMain's cleanup. The planner asked which intent the change set serves and how to clean up; no owner answered, and the arbiter decided (assumed): it serves milestone-25, and each package sweeps its own prefix's folders older than an hour at start through one shared helper, tested with a stale and a fresh folder. No per-run directory deleted by make test. Nothing is left open.

## title

Test temp folder cleanup, autonomous plan interview

## Links

### In

- derived-from [test-temp-cleanup-by-prefix](../decision/test-temp-cleanup-by-prefix.md)

## History

- `01M6M1M10N7AP6F0DHGANKXEA1` Fix Test temp folders leak and fill the disk: TestMain in internal/api/budget_test.go removes the generated rhiza-api-budget-* model (~146MB) only when the test binary exits normally, so every go test timeout (common under load) leaves one behind; quiettest leaves rhiza-quiet-* the same way, and slice/index budgets their own. On the VM 124 budget folders (~18GB) and 369 quiet folders piled up in TMPDIR by 2026-11-04 and the disk fell to 19GB free. Have each package remove its own stale folders (older than an hour, by prefix) when it starts, or put them under one per-run directory make test deletes, and test it
