# web-type-scale-qa-1

evidence · `01M7B9XGG73DH6D80Z3RTGJ8CD` · actor [claude-code-9fa50602](../actor/claude-code-9fa50602.md)

## result

PASS

## statement

Scratch build of HEAD run as a root daemon over a sandbox root of six projects (HOME and RHIZA_* isolated, web on 127.0.0.1:18931). Served rhiza-tokens.css holds exactly the 13 --fs and 11 --lh tokens; served rhiza.css has no font-size, line-height or font shorthand outside var() tokens (grep); go test ./internal/web ok.

## Links

### Out

- verifies [web-type-scale](../capability/rhiza/web-monitor/web-type-scale.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
