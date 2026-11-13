# web-project-needs-you-verified

evidence · `01M7B56M5YBK2ZP1MFVWVGERFN` · actor [claude-code-2d2790af](../actor/claude-code-2d2790af.md)

## result

PASS

## statement

At 57ef02502 re-ran the Checks block: go test -trimpath ./internal/web ./internal/api -run TestProjectNeedsYou (exit 0, TestProjectNeedsYou PASS), go test -trimpath ./internal/api -run Web (exit 0), rhiza ship generate (exit 0, 0 files changed), rhiza render (exit 0, no tracked change); go vet and gofmt clean. Diff read: ProjectView.NeedsYou (internal/web/projects.go) filled in WebProjectsSource (internal/api/web_root.go) from projectNeedsYou (internal/api/web_needs_you.go), one function that lists questions/open stops, sign-offs, failed runs and stuck runs (live run, agent step with a session, last journal entry over 20 minutes old). The test covers question+sign-off+failed+hold = 3, plus a 21-minute-quiet run = 4, a 19-minute run and a 5-minute run add nothing, and an empty project = 0. The acceptance's count of 5 with a hold is superseded by decision needs-you-count-excludes-machine-holds (count 4, holds in no count), which the code follows; the Home claim is dropped per needs-you-count-scoped-to-frame. No stubs or TODOs; nothing outside the leaf changed beyond the run's actor record and rendered files.

## Links

### Out

- verifies [web-project-needs-you](../capability/rhiza/web-monitor/web-project-needs-you.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
