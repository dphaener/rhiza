# web-frame-header-qa-1

evidence · `01M7B9XG4ZQT52DS4NG5P25MY0` · actor [claude-code-9fa50602](../actor/claude-code-9fa50602.md)

## result

PASS

## statement

Scratch build of HEAD run as a root daemon over a sandbox root of six projects (HOME and RHIZA_* isolated, web on 127.0.0.1:18931). /p/beta/, /p/beta/runs and /p/beta/roadmap each served wordmark, picker button (data-rz-picker=/picker, current project name), Home/Runs/Roadmap nav with aria-current on the current page and the live indicator; no badge when nothing waits, badge 1 on zeta with a failed run; no other project's name in the page. Served CSS: --header-h 52px, phone header min-height 56px with the nav on a 44px row. The merged /inbox and /projects pages have no picker button (finding recorded); the daemon pill shows on Home only.

## Links

### Out

- verifies [web-frame-header](../capability/rhiza/web-monitor/web-frame-header.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
