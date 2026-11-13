# web-page-head-spacing-qa-1

evidence · `01M7B9XG93KRSKZQKX4R8KC9ZA` · actor [claude-code-9fa50602](../actor/claude-code-9fa50602.md)

## result

PASS

## statement

Scratch build of HEAD run as a root daemon over a sandbox root of six projects (HOME and RHIZA_* isolated, web on 127.0.0.1:18931). Served rhiza.css: .rz-main top padding var(--space-6) (space-7 in the desktop query, per tests), .rz-pagehead > div gap var(--space-3), .rz-main gap var(--space-6), .rz-row bottom border var(--rule), tappable rows min-height 44px; go test ./internal/web ok.

## Links

### Out

- verifies [web-page-head-spacing](../capability/rhiza/web-monitor/web-page-head-spacing.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
