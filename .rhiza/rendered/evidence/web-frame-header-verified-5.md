# web-frame-header-verified-5

evidence · `01M7BE8W58M4B6TTN3G0QTA3D8` · actor [claude-code-1d022911](../actor/claude-code-1d022911.md)

## result

PASS

## statement

At 42739f9a7: acceptance go test -trimpath ./internal/web -run 'TestLayoutGolden|TestHeaderNoProjectList' exit 0; go test -trimpath ./internal/web/... exit 0 (incl. new TestHeaderDaemonPillOnEveryPage: pill on /, /runs, /roadmap, /p/a/*, /projects, /inbox); the new test fails to build on HEAD~1 (no web.WithDaemon), as the Checks block says; go test -trimpath -short -timeout 30m ./internal/api/... exit 0 (incl. TestWebRootHandlerServesEveryProject asserting the pill on root pages); go vet ./internal/web ./internal/api exit 0; rhiza ship generate and rhiza render exit 0 with no diff. The diff only touches the header's daemon pill (web.go render, WithDaemon, WebRootHandler) and its tests; no stubs or TODOs.

## Links

### Out

- verifies [web-frame-header](../capability/rhiza/web-monitor/web-frame-header.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
