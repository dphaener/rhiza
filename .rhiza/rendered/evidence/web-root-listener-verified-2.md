# web-root-listener-verified-2

evidence · `01M6SVHTC02M7EE0EMB9RX803G` · actor [claude-code-a8f11b1e](../actor/claude-code-a8f11b1e.md)

## result

PASS

## statement

Verifier claude-code-a8f11b1e at 970b7136 (attempt 2). Read the fix diff: WebRootHandler now copies ctx once (snap := *ctx) when it is built, synchronously in runDaemon before the poll loop, and the resolver, WebProjectsSource and WebMergedInboxSource derive every projectCtx from that copy. That resolves the confirmed quality finding f1 (the live ctx raced the poll's ensureRepoSession writes). runDaemon still builds a single web.Manager with WebRootHandler, so one listener serves the home project at / and each registered project at /p/<project>/. Ran `go test -trimpath -race -count=1 -run 'WebRoot|WebProject|Projects|MergedInbox|Daemon' ./internal/api/`: ok, no races, including TestWebRootHandlerServesEveryProject (two projects' pages answer through one handler). Ran `go test -trimpath -count=1 ./internal/web/...`: ok in a scratch copy outside ~/code. Under ~/code only TestOverviewReadsDaemonFromRepo fails, because it reads the machine's root daemon, which this change does not touch. go vet and gofmt are clean on internal/api and internal/web.

## Links

### Out

- verifies [web-root-listener](../capability/rhiza/web-root-listener.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
