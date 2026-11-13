# web-frame-header-verified-2

evidence · `01M7BAYJQSHP56C62PT1C3JRTT` · actor [claude-code-536e02f8](../actor/claude-code-536e02f8.md)

## result

PASS

## statement

At 51868432f (attempt 2, review cycle 2 fix): re-ran the Checks block. go test -trimpath ./internal/web -run TestHeaderPickerOnRootPages on HEAD~1 with the new test exits 1; go test -trimpath ./internal/web/... exits 0; the acceptance run go test -trimpath ./internal/web -run 'TestLayoutGolden|TestHeaderNoProjectList|TestHeaderPickerOnRootPages' passes (all three PASS); go vet ./internal/web exits 0; rhiza ship generate and rhiza render exit 0 with no tracked change. The diff only adds page.Root, sets it on the root's merged inbox and all-projects home (mounted only in root.go), and shows the picker button labelled All projects there; no TODOs or stubs, no change outside the leaf. Picker count, /picker link, nav aria-current, no other project names and no badge when nothing waits stay covered by the header tests.

## Links

### Out

- verifies [web-frame-header](../capability/rhiza/web-monitor/web-frame-header.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
