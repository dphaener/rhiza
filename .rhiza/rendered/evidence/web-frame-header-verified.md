# web-frame-header-verified

evidence · `01M7B9BY7CR6GMRNHBWC7K63GC` · actor [claude-code-0aff87f4](../actor/claude-code-0aff87f4.md)

## result

PASS

## statement

At f08934769: acceptance `go test -trimpath ./internal/web -run 'TestLayoutGolden|TestHeaderNoProjectList'` exit 0; TestHeaderNoBadgeWhenNothingWaits (failure case, count 0 shows no rz-count) passes; `go test -trimpath ./internal/web/...` ok; `rhiza ship generate` and `rhiza render` exit 0 with a clean tree. `go test -trimpath ./internal/api` hit go's default 10m package timeout twice on a loaded 4-core host (inside TestReadBudgetApply at 53s); rerun as the Makefile does (-timeout 20m class, -skip TestReadBudget) it passed (ok, 910s). Diff read: layout.html header has wordmark, picker button (data-rz-picker=/picker, current project name, needs-you badge only when >0, from HandlerWithNeedsYou fed by WebNeedsYou of that project alone), Home/Runs/Roadmap nav with aria-current, daemon pill, live indicator; Projects/Inbox links gone; phone header min-height 56px with the nav on a 44px row, desktop --header-h 52px; the test asserts a two-project root's /p/a/runs names neither /p/b nor b. Goldens regenerated. No TODOs, no changes outside the leaf.

## Links

### Out

- verifies [web-frame-header](../capability/rhiza/web-monitor/web-frame-header.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
