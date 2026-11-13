# web-project-needs-you-qa-2

evidence · `01M7BB8RSEA2445F8SWCY8YJYR` · actor [claude-code-9b577926](../actor/claude-code-9b577926.md)

## result

PASS

## statement

Scratch build of 4a9256c73 run as a root daemon over a sandbox root of six projects (HOME and RHIZA_* isolated, web on 127.0.0.1:18947). Real failed runs (command step run=false, run start --foreground): one in zeta, two in delta. Header badge and picker count were 1 and 2, every other project 0, All projects 3, matching rhiza inbox; questions added in a draft change set counted nowhere, as in rhiza inbox. TestProjectNeedsYou in internal/web and internal/api passes.

## Links

### Out

- verifies [web-project-needs-you](../capability/rhiza/web-monitor/web-project-needs-you.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
