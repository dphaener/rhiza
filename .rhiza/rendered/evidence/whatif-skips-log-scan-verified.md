# whatif-skips-log-scan-verified

evidence · `01M6BQZAQ02XK54JMWQM638CGW` · actor [claude-code-c8ac9d42](../actor/claude-code-c8ac9d42.md)

## result

PASS

## statement

Verified at 06744273. index.WhatIf now goes on from the stored as-is rows (model.WhatIfFrom) and decodes only the entries of cs, its in-flight closure and its landed dependencies (foldOne), not the whole log. Property: checkModel compares x.WhatIf against model.WhatIf over the full FoldLog (nodes, links, roots, problems) for every change set plus a missing one, and TestPropertyIncremental calls it. Its generator draws bases from any op, including unlanded ones, link.retire ops, and depends-on links whose ends include change set IDs. go test ./internal/index -run 'TestPropertyIncremental|TestWhatIf' -rapid.checks=1500 passed, and go test ./internal/index ./internal/model ./internal/changeset ./internal/api passed, including read_test.go's stored-versus-what-if comparison (go test ./internal/api -run TestRead). BenchmarkShowWhatIf on the 50k default model (gen.Default: 35000 nodes, 15000 links), -benchtime=8x -count=3: parent c403ceb0 1621/1625/1633 ms/op, this commit 244/240/245 ms/op, about 6.7x. make fmt-check and vet pass. make check's self-check fails only on the render of the verifier's own actor record, which was added after the commit.

## Links

### Out

- verifies [whatif-skips-log-scan](../capability/rhiza/views/fast-reads/whatif-skips-log-scan.md)

## History

- `01M6BH7ZQEKFG0H66J3X50B4TZ` A what-if read folds only the in-flight change sets it needs over the indexed as-is instead of the whole log
