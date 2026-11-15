# ci-build-cache-warmed-on-main

decision · `01M7FT14WNNJAFJJWBMDKXXM0W`

## reasoning

setup-go keys its cache on go.sum's hash, so it was saved once and never refreshed ('Cache hit occurred on the primary key ..., not saving cache'), and every job compiles the repo's packages again: the api 1/2 step is 149 s for 86 s of tests, budgets 125 s for 73 s, scripts 3/4 140 s for 112 s. A pull request can restore caches its base branch saved, and saving from main alone keeps the cache within the repo's quota. It adds no job to a pull request's run, which matters: run 43721082089 waited 3.5 min for a runner.

## statement

CI's Go build and module caches come from main: a workflow on push to main (and workflow_dispatch) builds every package's tests the way the jobs do (-trimpath, go test -run '^$' ./... and go vet ./...) and saves ~/.cache/go-build and ~/go/pkg/mod with actions/cache under a key naming the OS, Go version and commit; every Go job in ci.yml restores the newest such key by prefix and saves nothing, with setup-go's own cache off

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [fast-suite-finish-research](../source/fast-suite-finish-research.md)

### In

- decided-by [ci-fresh-build-cache](../capability/ci-fresh-build-cache.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
