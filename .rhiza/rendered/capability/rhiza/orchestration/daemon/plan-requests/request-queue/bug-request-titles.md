# bug-request-titles

capability · `01M6GDV6ERT7J1R6WGN7VH2SC1` · parent [request-queue](../request-queue.md)

## acceptance

TestBugRequestTitlesAreDistinct (internal/api): two bug requests' intents give display titles that are not "Fix" and have distinct slugs; plan-requests.txtar shows the daemon starting run start --new "Fix Login fails"; describe-golden.txtar regenerated

## statement

Each bug request's change set has its own title, so bug plan runs at once never share a branch or worktree

## History

- `01M6GDJW3R8PXW7QBAZVH0QS4W` Bug requests' change sets get distinct titles: a bug's plan run opens its change set with the intent 'Fix <text>' rather than 'Fix: <text>', so its display title, branch and worktree name the bug instead of all being 'Fix'
