# web-spacing-rules

capability · `01M79Q90R4MS5J2BXQ7DGXA1TS` · parent [web-monitor](../web-monitor.md)

## acceptance

go test -trimpath ./internal/web -run 'TestPartialsGolden|TestSpacingRules' passes: the goldens are regenerated, and TestSpacingRules parses static/css/rhiza.css and fails unless the page head's gap between kicker and title is var(--space-3) (12px), the panel head's gap between kicker and title is var(--space-1) (4px), the panel head's margin to the content is var(--space-4) (16px), the gap between panels is var(--space-5) (24px), the gap between page sections is var(--space-6) (32px), and panel padding is var(--space-4) (16px) by default and var(--space-5) (24px) inside the desktop media query; it also fails on any heading rule with a margin and any template with an inline margin style (failure case: a stylesheet with the panel gap at 16px fails).

## statement

Every web page follows the spacing rules in pages/frame-and-spacing.html: kicker to page title 12px, panel kicker to panel title 4px, panel title to content 16px, panels 24px apart, page sections 32px, panel padding 16px on a phone and 24px on desktop, gaps set on the parent and headings with no margin of their own.

## Links

### Out

- decided-by [web-spacing-scale](../../../decision/web-spacing-scale.md)

### In

- verifies [web-spacing-rules-verified](../../../evidence/web-spacing-rules-verified.md)
- verifies [web-spacing-rules-qa-1](../../../evidence/web-spacing-rules-qa-1.md)
- verifies [web-spacing-rules-qa-2](../../../evidence/web-spacing-rules-qa-2.md)
- verifies [pr-571-the-web-ui-s-frame](../../../evidence/pr-571-the-web-ui-s-frame.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
