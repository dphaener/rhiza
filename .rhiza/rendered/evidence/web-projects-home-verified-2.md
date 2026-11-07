# web-projects-home-verified-2

evidence · `01M6SW4VHK3577Y8JAXW8181FQ` · actor [claude-code-fa0195ce](../actor/claude-code-fa0195ce.md)

## result

PASS

## statement

At b70cca66 (attempt 2), from a scratch worktree outside ~/code, go test -trimpath -count=1 ./internal/web/... passed. That includes TestProjectsHomeReadsTheRoot: a root of two initialised projects (alpha, beta) and one bare one, rendered on the projects home, checking each row's live-run count, inbox count, health pill (matched against each project's own health report) and /p/<name>/ link, the bare one marked not initialised with no counts, and every page's nav linking to /projects. It also includes the new TestRootRedirectsUnscopedToHome: / serves the projects home, unscoped pages 302 to /p/<home>/..., and /static, /api, /events and the webmanifest stay served by home. go test ./internal/api -run TestWeb passed, including TestWebRootHandlerServesEveryProject, where / lists both projects and /runs redirects to /p/alpha/runs. I read the diff: it makes / the all-projects home, which resolves the attempt-1 compliance finding f1. Inside ~/code, TestOverviewReadsDaemonFromRepo fails only because it sees the live root daemon (a known environment issue); it passes outside ~/code.

## Links

### Out

- verifies [web-projects-home](../capability/rhiza/web-projects-home.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
