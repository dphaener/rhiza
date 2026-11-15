# link-endpoint-kinds-fast-verified

evidence · `01M7FW0S688E94D6W3WBX9AWQ5` · actor [claude-code-8097663b](../actor/claude-code-8097663b.md)

## result

PASS

## statement

At e12b6202d, under the machine lock with RHIZA_*/SHIP_*/REX_* unset: go test -trimpath -p 2 -count=1 -v ./internal/api -run '^TestLinkEndpointKinds$' exit 0, 15 parallel per-link-type subtests PASS, test 9.05 s, 10.33 s wall (under 15 s; was 81 s). go vet ./internal/api exit 0. The diff adds an atomic count and fails unless checked == len(op.LinkTypes()) x len(check.Ends())^2 (the pass confirms equality). Scratch mutation, not committed and reverted with git checkout: flipping the check.Allowed answer at its call in internal/api/write_link.go:102 for supersedes change set -> change set made the test fail with 'supersedes from change set to change set is not allowed, but got <nil>', naming the pair. (Flipping check.Allowed itself would change the test's oracle too, so the mutation goes at the call site, as the builder noted.) The code change touches internal/api/write_test.go only; the rest is the builder's staged actor record and rendered index. No stubs or TODOs.

## Links

### Out

- verifies [link-endpoint-kinds-fast](../capability/link-endpoint-kinds-fast.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
