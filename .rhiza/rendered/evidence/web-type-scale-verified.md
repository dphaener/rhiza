# web-type-scale-verified

evidence · `01M7B5A6ZXAF5F2XSTCWDP5ENP` · actor [claude-code-d240aae7](../actor/claude-code-d240aae7.md)

## result

PASS

## statement

At f41bce522 (attempt 2): go test -trimpath ./internal/web -run TestTypeScale exit 0 and go test -trimpath ./internal/web/... exit 0, as in the Checks: block. Attempt 1's finding is fixed: .rz-step::before is back to left: 10.5px, and across both leaf commits (980db15cb, f41bce522) rhiza.css and templates/change.html change only font, font-size and line-height values, all now var(--fs-*)/var(--lh-*); rhiza-tokens.css adds exactly the 13 size and 11 line-height tokens. The listed moves are in place (10.5px to 11px, 11.5px to 12px, 12.5px to 12px, 13.5/20 to 14/20, stat 28/32 to 28/34, phone titles to 32/36 and 28/34). The test reads pages/*.html. In a scratch copy, a 28px/32px stat value and a 10.5px row head in rhiza.css each make TestTypeScale fail. No TODOs, no new dependencies.

## Links

### Out

- verifies [web-type-scale](../capability/rhiza/web-monitor/web-type-scale.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
