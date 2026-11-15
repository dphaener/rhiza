# link-matrix-parallel-subtests

decision · `01M7FT149DPXVR1A9MMQVV9KRD`

## reasoning

It took 81 s serially, the single longest internal/api test and the long pole of CI's api 1/2 shard, because one loop copies the linkEndpoints fixture for each of about 1,500 cases. A refused link writes nothing, so its cases can share a copy (if one is wrongly allowed the test fails anyway); per-type subtests spread the rest over the cores. No case is dropped, so no coverage is lost.

## statement

TestLinkEndpointKinds keeps its full matrix (every link type x every end x every end) but runs one parallel subtest per link type, each with one fixture copy shared by the links it expects refused and a fresh copy only for each link it expects allowed

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [fast-suite-finish-research](../source/fast-suite-finish-research.md)

### In

- decided-by [link-endpoint-kinds-fast](../capability/link-endpoint-kinds-fast.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
