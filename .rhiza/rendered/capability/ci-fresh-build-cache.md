# ci-fresh-build-cache

capability · `01M7FT14YKGTQ47MDNF0MCX1YB`

## acceptance

go test -trimpath -p 2 -count=1 ./internal/ciworkflow passes, with tests pinning the new workflow (push to main and workflow_dispatch; one job that builds with -trimpath, go test -run '^$' ./... and go vet ./..., then saves ~/.cache/go-build and ~/go/pkg/mod with actions/cache under a key with runner.os, the Go version and github.sha) and pinning in ci.yml that setup-go's cache is off and every Go job (model, lint, scripts, api, rest, budgets) restores that cache by key prefix with restore-keys and saves none; the existing ciworkflow tests (jobs, matrices, check aggregation) still pass. The action versions match the ones ci.yml already uses where they overlap. Measure (orchestrator, after merge and one push to main): in the next code PR's CI run, gh run view <id> --log shows each Go job restoring a cache saved from main, and the api jobs' Run api step exceeds its go test time by under 20 s (2026-11-15: 63 s). Touches .github/workflows/ci.yml, .github/workflows/go-cache.yml (new) and internal/ciworkflow/ciworkflow_test.go.

## statement

CI's Go jobs start from a build cache of main's latest commit: a workflow on push to main builds every package's tests with -trimpath and saves the Go build and module caches under a key naming the OS, Go version and commit, and each Go job in ci.yml restores the newest such cache by prefix instead of setup-go's go.sum-keyed one, so the jobs compile only what a pull request changed.

## Links

### Out

- decided-by [ci-build-cache-warmed-on-main](../decision/ci-build-cache-warmed-on-main.md)

### In

- verifies [ci-fresh-build-cache-verified](../evidence/ci-fresh-build-cache-verified.md)
- verifies [ci-fresh-build-cache-qa-1](../evidence/ci-fresh-build-cache-qa-1.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
