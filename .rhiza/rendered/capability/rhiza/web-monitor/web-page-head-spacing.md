# web-page-head-spacing

capability · `01M7B25EAACMG725WJG09RAG5E` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run 'TestPageHeadSpacing' passes: it parses static/css/rhiza.css and fails unless the space above the page head (.rz-main's top padding) is var(--space-6) by default and var(--space-7) inside the desktop media query, the gap between the page title and its line is var(--space-3), the page head's margin to the content is var(--space-6), a list row (.rz-row) has 10px padding above and below and a var(--rule) bottom border, and a tappable row (a row that is or holds a link) has min-height 44px (failure cases: a 48px top padding outside the desktop query fails, and a tappable row with no min-height fails).

## statement

Every web page follows the page-head and row rules in pages/frame-and-spacing.html: header to page head space-6 (32px) on a phone and space-7 (48px) on desktop, page title to its line space-3 (12px), page head to content space-6 (32px), and list rows 10px above and below with a 1px line between and a 44px minimum height when the row is tappable.

## Links

### Out

- decided-by [web-spacing-scale](../../../decision/web-spacing-scale.md)

### In

- verifies [web-page-head-spacing-verified](../../../evidence/web-page-head-spacing-verified.md)
- verifies [web-page-head-spacing-qa-1](../../../evidence/web-page-head-spacing-qa-1.md)
- verifies [web-page-head-spacing-qa-2](../../../evidence/web-page-head-spacing-qa-2.md)
- verifies [pr-571-the-web-ui-s-frame](../../../evidence/pr-571-the-web-ui-s-frame.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
