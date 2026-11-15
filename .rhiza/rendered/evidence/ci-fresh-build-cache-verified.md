# ci-fresh-build-cache-verified

evidence · `01M7G2QRTD3897QHQV6C6DN118` · actor [claude-code-d2313fd1](../actor/claude-code-d2313fd1.md)

## result

PASS

## statement

At 6dd79ebd6: go test -trimpath -p 2 -count=1 ./internal/ciworkflow exit 0 (all 11 tests pass, incl. new TestCIGoJobsRestoreMainsCacheAndSaveNone and TestGoCacheWarmsFromMain and the existing jobs/matrices/check-aggregation tests); rhiza ship generate exit 0 (0 changed); rhiza render exit 0. Diff read: go-cache.yml triggers push to main and workflow_dispatch, one job builds with go build -trimpath, go test -trimpath -run '^$' ./... and go vet ./..., saves ~/.cache/go-build and ~/go/pkg/mod with actions/cache@v6 under go-runner.os-go-version-github.sha. ci.yml: all 6 Go jobs (model, lint, scripts, api, rest, budgets) set setup-go cache: false and restore via actions/cache/restore@v6 with the prefix restore-keys, none saves; tests pin both. checkout@v7 and setup-go@v7 match ci.yml's existing versions. Only the three named files plus log/actor record changed; no TODOs. Measure lines left to the orchestrator.

## Links

### Out

- verifies [ci-fresh-build-cache](../capability/ci-fresh-build-cache.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
