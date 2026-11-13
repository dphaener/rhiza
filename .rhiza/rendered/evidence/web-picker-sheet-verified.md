# web-picker-sheet-verified

evidence · `01M7B65CPXPKX0HSQNJ2SA2ZYR` · actor [claude-code-3da9d12d](../actor/claude-code-3da9d12d.md)

## result

PASS

## statement

At ebedabfa3 re-ran the Checks block: go test -trimpath ./internal/web -run 'TestProjectPicker|TestSwitcherNamesReadOnlyForPages' -v (both PASS, exit 0), go test -trimpath ./internal/web/... (ok), go test -trimpath ./internal/api -run Web (ok), rhiza ship generate (exit 0, 0 changed), rhiza render (exit 0, tracked tree unchanged). Read the diff: picker.go serves /picker from the ProjectsSource once per request, listing initialised projects with live state (N running / idle) and needs-you counts, sorted with items first, links keeping the page kind via pageKind (detail pages keep their list), an All projects row linking /inbox, and search only from six projects; TestProjectPicker asserts each of these (2 vs 6 initialised projects for search). RootWithSwitcher, ProjectNames, WebProjectNames, scope.Switcher and the inline switcher template are removed with no remaining references; TestSwitcherNamesReadOnlyForPages now asserts assets, fragments and page renders under /p/a/ never read the source and one GET /picker reads it once. CSS uses defined tokens (bottom sheet by default, menu from 721px); no TODOs or stubs; tests use fakes only.

## Links

### Out

- verifies [web-picker-sheet](../capability/rhiza/web-monitor/web-picker-sheet.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
