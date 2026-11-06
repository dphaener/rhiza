# web-overview-live-qa-1

evidence · `01M6Q95MP300MR4P9E8K0FCBHP` · actor [claude-code-885ce195](../actor/claude-code-885ce195.md)

## result

PASS

## statement

Same live daemon: / carries data-rz-fragment and data-rz-events on each section (daemon: daemon inbox change; inbox: inbox run change; runs: run daemon; milestones: change). GET /fragments/daemon, /runs, /inbox and /milestones each returned 200 with only that section's panel, no <html> and no rz-section wrapper. /events sent ': ready' and keep-alives, and 'event: run' with the run's id when its journal grew; the runs fragment then showed the new step.

## Links

### Out

- verifies [web-overview-live](../capability/rhiza/web-monitor/web-overview-live.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
