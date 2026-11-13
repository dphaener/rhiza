# web-project-needs-you-qa-1

evidence · `01M7B9XFTQCRP7MXDFQF19NR5H` · actor [claude-code-9fa50602](../actor/claude-code-9fa50602.md)

## result

PASS

## statement

Scratch build of HEAD run as a root daemon over a sandbox root of six projects (HOME and RHIZA_* isolated, web on 127.0.0.1:18931). With only the daemon's machine-wide load-high hold active every project's count was 0 in the header and /picker (hold not counted). A real failed run in zeta (workflow with a failing command step, run start --foreground) made zeta's header badge and picker count 1, the other five stayed 0, matching rhiza inbox in zeta.

## Links

### Out

- verifies [web-project-needs-you](../capability/rhiza/web-monitor/web-project-needs-you.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
