# m40-no-systemd-scopes

decision · `01M7NVQ8ER86FPA346X7PYA6W3`

## question

Host limits: how should the daemon enforce them for any workflow? (scopes plus test slots, session scopes only, or test slots only)

## quote

A systemd scope seems like overkill to me.

## quote_where

owner, 2026-11-17 orchestrator chat, answer to question 1

## statement

The daemon enforces no host limits through systemd scopes or cgroups (Linux-only, and Rhiza also runs on macOS) and keeps no host test-slot count (the daemon cannot tell which of a repo's commands are tests; a repo queues its own tests, as this repo's make test does with its machine-wide lock)

## Links

### Out

- derived-from [m40-interview](../source/m40-interview.md)

### In

- decided-by [milestone-40](../intent/milestone-40.md)

## History

- `01M7MXN3D49WSG4W8YEQEF4WFJ` The daemon runs any workflow, built-in or bring-your-own, from start to finish with no orchestrator beside it: it enforces host limits on child sessions and their tests, detects and recovers stalled steps, keeps the box healthy before and after runs, retries failed runs from the right step with their findings, hands planners the notes queued with a request, and records every outside intervention so intervention stays the exception. Brief with the owner's words: /home/rhiza/orchestrator/hands-off-milestone-interview.md
