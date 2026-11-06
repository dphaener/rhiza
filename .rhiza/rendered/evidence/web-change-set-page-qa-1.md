# web-change-set-page-qa-1

evidence · `01M6Q7HGF2KTHCVZ9EVNG8Z9PY` · actor [claude-code-3eb24189](../actor/claude-code-3eb24189.md)

## result

FAIL

## statement

At 75d6f69b, on the same running daemon, GET /changes/01M6PFH625P3BYQAB8P6J38BBT (this change set) returned 200 with its title, intent, kind, accepted pill, what it serves and a Proposed/Accepted history, but with '0 leaves, 0 with evidence' and 'Depends on —'. Its change.show lists two leaves, and the roadmap gives it two depends_on. /api/show?ref=<leaf> with no view came back 404 in the as-is view, because the clone had no current change set. Once a change set was current, the default view became to-be and the same page showed 2 leaves, 2 with evidence. Depends on stayed empty for every change set I opened, because the page reads only depends-on links that the change set's own ops add, and here the milestone plan added them. After I copied this run's journal into the clone, Runs showed the ship-change run linking to /runs/<id>, with a blank time. Unknown ULID and non-ID refs gave 404.

## Links

### Out

- verifies [web-change-set-page](../capability/rhiza/web-monitor/web-change-set-page.md)

## History

- `01M6PFH625P3BYQAB8P6J38BBT` Web Roadmap pages: milestone groups and the backlog, and a change set's page with its intent, state, leaves, questions, runs and history
