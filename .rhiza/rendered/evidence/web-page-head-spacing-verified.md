# web-page-head-spacing-verified

evidence · `01M7B7HVN80W97BF0ZCXDF40Q1` · actor [claude-code-110c0e8b](../actor/claude-code-110c0e8b.md)

## result

PASS

## statement

At cd6321dfe (attempt 2): re-ran the Checks block, all exit 0: go test -trimpath ./internal/web -run 'TestPageHeadSpacing' (PASS), go test -trimpath ./internal/web (ok), rhiza ship generate, rhiza render. rhiza.css: .rz-main padding top var(--space-6) by default and var(--space-7) under (min-width: 721px); .rz-pagehead > div gap var(--space-3); .rz-main gap var(--space-6); .rz-row padding 10px by default and 10px under (max-width: 720px) (was 12px, the attempt-1 gap); .rz-row border-bottom var(--rule) solid var(--line); .rz-row:has(a), a.rz-row min-height 44px. Failure cases tried on a scratch worktree copy: .rz-main top padding 48px outside the desktop query fails the test (got 48px, want var(--space-6)); removing the tappable-row min-height rule fails it; phone row padding 12px fails it. Diff touches only internal/web/spacing_test.go, internal/web/static/css/rhiza.css and .rhiza actor log/rendered files; no stubs or TODOs.

## Links

### Out

- verifies [web-page-head-spacing](../capability/rhiza/web-monitor/web-page-head-spacing.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
