# web-spacing-rules-qa-1

evidence · `01M7B9XGD9JY8EP0W2BJ8HPH4T` · actor [claude-code-9fa50602](../actor/claude-code-9fa50602.md)

## result

PASS

## statement

Scratch build of HEAD run as a root daemon over a sandbox root of six projects (HOME and RHIZA_* isolated, web on 127.0.0.1:18931). Served rhiza.css: page head kicker-title gap space-3, panel head gap space-1, panel head margin-bottom space-4, .rz-stack gap space-5, .rz-main gap space-6; templates carry no inline margin style (grep); go test ./internal/web ok.

## Links

### Out

- verifies [web-spacing-rules](../capability/rhiza/web-monitor/web-spacing-rules.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
