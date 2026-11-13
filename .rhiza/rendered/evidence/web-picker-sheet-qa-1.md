# web-picker-sheet-qa-1

evidence · `01M7B9XG0GY7H3591GC8FKPEG8` · actor [claude-code-9fa50602](../actor/claude-code-9fa50602.md)

## result

PASS

## statement

Scratch build of HEAD run as a root daemon over a sandbox root of six projects (HOME and RHIZA_* isolated, web on 127.0.0.1:18931). GET /picker?path=/p/beta/roadmap listed all six projects with idle state and counts, beta aria-current, links kept the page kind (/p/x/roadmap; /p/beta/runs/<id> gave /p/x/runs; no path gave /p/x/), All projects linked /inbox, search field present at six projects; with zeta's failed run zeta sorted first and All showed 1. Unknown project or bad path query render cleanly (200, nothing current); POST /picker 404. Project pages' markup names no other project.

## Links

### Out

- verifies [web-picker-sheet](../capability/rhiza/web-monitor/web-picker-sheet.md)

## History

- `01M79NKERTK62YZM9C42BJ0JW3` The web UI's frame: one header and navigation for every screen with a real project picker (search, each project's state and needs-me count) in place of the header's project list, and shared components on one spacing and type scale
