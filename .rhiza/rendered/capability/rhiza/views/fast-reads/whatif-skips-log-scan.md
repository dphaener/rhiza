# whatif-skips-log-scan

capability · `01M6BH9DH9XCPBR849KJ72RSYB` · parent [fast-reads](../fast-reads.md)

## acceptance

A property test over generated logs (with change sets based on unlanded ones, retires and links to change sets) finds the new what-if equal to model.WhatIf over the full fold, model and problems both; read_test.go's stored-versus-what-if comparison still passes; BenchmarkShowWhatIf on the 50k model improves at least fivefold, recorded in the Evidence

## statement

A what-if read decodes only the in-flight change sets it needs and the landed ones their revisions name, over the index's stored as-is rows, instead of folding the whole log, and gives the same model and problems as the full fold

## Links

### Out

- decided-by [whatif-folds-inflight-over-asis](../../../../decision/whatif-folds-inflight-over-asis.md)

### In

- verifies [whatif-skips-log-scan-verified](../../../../evidence/whatif-skips-log-scan-verified.md)
- verifies [whatif-skips-log-scan-qa-1](../../../../evidence/whatif-skips-log-scan-qa-1.md)
- verifies [pr-173-what-if-skips-the-log-scan](../../../../evidence/pr-173-what-if-skips-the-log-scan.md)

## History

- `01M6BH7ZQEKFG0H66J3X50B4TZ` A what-if read folds only the in-flight change sets it needs over the indexed as-is instead of the whole log
