# web-change-set-page

capability · `01M6PFKANQ23RTM5QSR6H1CGBX` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test renders /changes/<id> for a fixture change set with leaves (one with Evidence), a blocking question, a dependency and a run, and asserts each part and its links, and a 404 for an unknown change set.

## statement

A change set's page at /changes/<id> shows its title, intent, kind and state, what it serves and depends on, its leaves with each statement, acceptance and whether it has Evidence, its open questions, its runs linking to their run pages, and its lifecycle history.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [web-change-set-page-verified](../../../evidence/web-change-set-page-verified.md)
- verifies [web-change-set-page-qa-1](../../../evidence/web-change-set-page-qa-1.md)
- verifies [web-change-set-page-verified-2](../../../evidence/web-change-set-page-verified-2.md)
- verifies [web-change-set-page-qa-2](../../../evidence/web-change-set-page-qa-2.md)
- verifies [web-change-set-page-verified-5](../../../evidence/web-change-set-page-verified-5.md)
- verifies [web-change-set-page-qa-3](../../../evidence/web-change-set-page-qa-3.md)
- verifies [pr-319-web-ui-roadmap-pages](../../../evidence/pr-319-web-ui-roadmap-pages.md)

## History

- `01M6PFH625P3BYQAB8P6J38BBT` Web Roadmap pages: milestone groups and the backlog, and a change set's page with its intent, state, leaves, questions, runs and history
