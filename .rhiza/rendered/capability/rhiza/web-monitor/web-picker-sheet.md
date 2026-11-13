# web-picker-sheet

capability · `01M7B25DYN1X19CGG3X73K97J5` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run 'TestProjectPicker|TestSwitcherNamesReadOnlyForPages' passes: TestProjectPicker GETs /picker for a page of a two-project root and finds both listed with live state and needs-you counts, the one with items first, each link keeping the page kind (from /p/a/roadmap, b links to /p/b/roadmap), All projects linking to the merged inbox, and the search field present only with six or more projects; TestSwitcherNamesReadOnlyForPages, rewritten, fails if an asset, a fragment or a page render under /p/<project>/ calls the ProjectsSource, and if one GET /picker does not call it exactly once.

## statement

The header's project button opens a picker (a bottom sheet on a phone, a menu on desktop, as in pages/project-picker.html) listing every registered project with whether it is live and its needs-you count (both from web-project-needs-you), projects with something waiting first, plus an All projects entry for the merged inbox; switching keeps the same kind of page, and a search field appears once there are more than five projects. The list is not in the page: the root serves it at /picker, given the current page's path, and reads the ProjectsSource only for that request, so assets, fragments and page renders still never read it. This leaf retires RootWithSwitcher's ProjectNames list and the inline project-name switcher in templates/layout.html; web-frame-header, built after it, rewrites the header around the picker button.

## Links

### Out

- depends-on [web-project-needs-you](web-project-needs-you.md)
- decided-by [web-project-picker](../../../decision/web-project-picker.md)

### In

- depends-on [web-frame-header](web-frame-header.md)
- verifies [web-picker-sheet-verified](../../../evidence/web-picker-sheet-verified.md)
- verifies [web-picker-sheet-qa-1](../../../evidence/web-picker-sheet-qa-1.md)
- verifies [web-picker-sheet-qa-2](../../../evidence/web-picker-sheet-qa-2.md)
- verifies [pr-571-the-web-ui-s-frame](../../../evidence/pr-571-the-web-ui-s-frame.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
