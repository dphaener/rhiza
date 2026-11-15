# link-endpoint-kinds-fast

capability · `01M7FT14BJ76NCG2N32YX9190C`

## acceptance

go test -trimpath -p 2 -count=1 -v ./internal/api -run '^TestLinkEndpointKinds$' passes in under 15 s of wall time on the VM (2026-11-15: 81 s), and the test fails unless the number of cases it checked equals len(op.LinkTypes()) x len(check.Ends())^2. Failure case (a scratch mutation, not committed): flipping check.Allowed's answer for one pair makes the test fail naming that pair. Touches internal/api/write_test.go only.

## statement

TestLinkEndpointKinds checks every (link type, from end, to end) of op.LinkTypes() x check.Ends() x check.Ends() as before, in parallel subtests per link type that copy the fixture only for links they expect allowed, so it takes seconds instead of 81 s.

## Links

### Out

- decided-by [link-matrix-parallel-subtests](../decision/link-matrix-parallel-subtests.md)

### In

- verifies [link-endpoint-kinds-fast-verified](../evidence/link-endpoint-kinds-fast-verified.md)
- verifies [link-endpoint-kinds-fast-qa-1](../evidence/link-endpoint-kinds-fast-qa-1.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
