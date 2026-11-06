# web-change-set-page-qa-2

evidence · `01M6Q91VN2VZ65B931R6VSXKET` · actor [claude-code-0ebbd179](../actor/claude-code-0ebbd179.md)

## result

FAIL

## statement

At 1394f316, on the same daemon with no current change set (rhiza show <leaf> in the clone returns not_found), GET /changes/01M6PFH625P3BYQAB8P6J38BBT returned 200. It showed the title, intent, kind, accepted pill and what the set serves; both leaves with statement, acceptance and 'has evidence' ('2 leaves, 2 with evidence'); Depends on 'Web UI server' and 'Web UI design system', both landed and linking to /changes/<id>; the ship-change run linking to /runs/<id> with time 2026-11-06 00:57; and the Proposed and Accepted history. I added a question blocking web-roadmap-page with --change, still with no current change set, and the page listed it under Open questions. Unknown ULID and non-ID refs return 404, and POST returns 405. So cycle 1's f1, f2 and f4 are fixed for unlanded change sets. But on landed change sets, for example 01M6J965 (change.show lists 3 leaves and 6 created questions) and 01M6FENG (1 leaf), the page shows '0 leaves, 0 with evidence', Serves —, and no questions. /api/show?view=what-if:<landed> answers 409 'is landed ... so it has no what-if view', while the default view finds the leaf.

## Links

### Out

- verifies [web-change-set-page](../capability/rhiza/web-monitor/web-change-set-page.md)

## History

- `01M6PFH625P3BYQAB8P6J38BBT` Web Roadmap pages: milestone groups and the backlog, and a change set's page with its intent, state, leaves, questions, runs and history
