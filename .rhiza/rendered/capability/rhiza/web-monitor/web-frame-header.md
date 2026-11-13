# web-frame-header

capability · `01M79Q90M4F549HKEWB54K9B9P` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run 'TestLayoutGolden|TestHeaderNoProjectList' passes with the goldens regenerated from the new header: the picker button carries the needs-you count and points at /picker, the nav names Home, Runs and Roadmap with aria-current on the current page, and a two-project root's page names no project but the current one anywhere in its body (failure case: a project with nothing waiting shows no count badge).

## statement

Every web page shares one header as in pages/frame-and-spacing.html and pages/home.html: the Rhiza wordmark, the project picker button with the current project's needs-you count (the count web-project-needs-you defines, read from that project alone), the Home, Runs and Roadmap nav, the daemon pill and the live indicator; 56px tall on a phone with the nav on its own 44px row, 52px on desktop; no list of project names in the header or anywhere in the page's markup, the picker's list being loaded from /picker (web-picker-sheet) when the button opens it. It is built after web-project-needs-you and web-picker-sheet, on the layout web-picker-sheet leaves without the inline switcher.

## Links

### Out

- depends-on [web-project-needs-you](web-project-needs-you.md)
- depends-on [web-picker-sheet](web-picker-sheet.md)
- decided-by [web-project-picker](../../../decision/web-project-picker.md)

### In

- verifies [web-frame-header-verified](../../../evidence/web-frame-header-verified.md)
- verifies [web-frame-header-qa-1](../../../evidence/web-frame-header-qa-1.md)
- verifies [web-frame-header-verified-2](../../../evidence/web-frame-header-verified-2.md)
- verifies [web-frame-header-qa-2](../../../evidence/web-frame-header-qa-2.md)
- verifies [web-frame-header-verified-5](../../../evidence/web-frame-header-verified-5.md)
- verifies [pr-571-the-web-ui-s-frame](../../../evidence/pr-571-the-web-ui-s-frame.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
