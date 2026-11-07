# web-root-listener-verified

evidence · `01M6SPH3WYZSXSYYT4635ZC7A7` · actor [claude-code-b04195b2](../actor/claude-code-b04195b2.md)

## result

PASS

## statement

At de2323a2 the verifier read the diff: runDaemon builds one web.Manager whose handler is WebRootHandler, which serves the home project at / and each registered, initialised project at /p/<project>/ through web.Root. It resolves projects from the root's rhiza.toml on each request, gives 404 for unknown or uninitialised projects and redirects /p/<name> to /p/<name>/. No other listener is opened (the only web.Manager/webFactory use is in daemon.go). go test -trimpath ./internal/api -run TestWebRootHandlerServesEveryProject PASS: a root of two initialised projects (alpha, beta) plus an uninitialised one answers /p/alpha and /p/beta /api/roadmap and /runs with 200, each /api/status from its own model, home /runs 200, bare and unknown 404. go test -trimpath ./internal/api -run 'TestWebRoot|TestDaemon|TestWeb' PASS. go test -trimpath ./internal/web/... PASS (TestRootRoutesProjects included) when run from a normal worktree at this commit. Run from the run's checkout under .git/rhiza/checkouts, TestOverviewReadsDaemonFromRepo fails for an environmental reason, not because of this change. go vet and gofmt are clean. Links in the pages are still absolute (/runs, /roadmap); scoping them is the web-project-switcher leaf's job.

## Links

### Out

- verifies [web-root-listener](../capability/rhiza/web-root-listener.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
