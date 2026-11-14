# web-frame-header-qa-2

evidence · `01M7BB8RM4G3P1KA1D8VQH7V54` · actor [claude-code-9b577926](../actor/claude-code-9b577926.md)

## result

PASS

## statement

Scratch build of 4a9256c73 run as a root daemon over a sandbox root of six projects (HOME and RHIZA_* isolated, web on 127.0.0.1:18947). /p/zeta/, /p/delta/runs, /p/beta/roadmap: wordmark, picker button (data-rz-picker=/picker, current project), Home/Runs/Roadmap nav, live indicator; badge 1 on zeta and 2 on delta (failed runs), no badge on beta; /p/zeta/runs names no other project. Cycle 2 fix: /, /projects and /inbox now carry the picker button labelled All projects, and /picker?path=/inbox lists every project with All projects linking /inbox. Scoped tests (TestLayoutGolden, TestHeader*) pass.

## Links

### Out

- verifies [web-frame-header](../capability/rhiza/web-monitor/web-frame-header.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
