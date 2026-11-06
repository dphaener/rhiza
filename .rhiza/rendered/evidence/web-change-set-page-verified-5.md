# web-change-set-page-verified-5

evidence · `01M6Q9AZGAJHFGN4A787ZV8VBE` · actor [claude-code-49cfd641](../actor/claude-code-49cfd641.md)

## result

PASS

## statement

At 95377cde ran go test -trimpath -count=1 ./internal/web/... (ok) and go vet ./internal/web (clean). TestChangeSetPage builds a fixture with no current change set: a proposed change set with two leaves (one verified by Evidence), a question blocking a leaf, a depends-on link added by another landed change set to a landed dependency, and a run. It asserts /changes/<id> shows title, intent kind, state, each leaf's statement and acceptance, has evidence / no evidence yet, '2 leaves, 1 with evidence', the open question, the dependency's title, link and state, the run linking to /runs/runone with its start time 2026-10-24 11:00, and the Proposed history entry, with no handles or IDs shown. It also renders the landed dependency's own page (its leaf, acceptance and '1 leaf, 0 with evidence'), and asserts 404 for an unknown ID and an unknown ref and 405 for POST. Read the diff: showRef reads what-if:<cs> for draft/proposed/accepted and as-is otherwise, dependencies come from the roadmap operation, and run started_at is tagged, so the earlier attempts' findings are addressed.

## Links

### Out

- verifies [web-change-set-page](../capability/rhiza/web-monitor/web-change-set-page.md)

## History

- `01M6PFH625P3BYQAB8P6J38BBT` Web Roadmap pages: milestone groups and the backlog, and a change set's page with its intent, state, leaves, questions, runs and history
