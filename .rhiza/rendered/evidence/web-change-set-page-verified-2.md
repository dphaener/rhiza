# web-change-set-page-verified-2

evidence · `01M6Q8KQ5K0ABKW9DAJ6FQJ3GZ` · actor [claude-code-31bbabd7](../actor/claude-code-31bbabd7.md)

## result

PASS

## statement

At ef7293b5 ran go test -trimpath ./internal/web/... (ok), go vet and gofmt -l on internal/web (clean). TestChangeSetPage renders /changes/<id> for a fixture change set with two leaves (one with Evidence: 'has evidence' / 'no evidence yet', '2 leaves, 1 with evidence'), statements and acceptance, a blocking open question, a dependency from a landed plan linking to /changes/<base> with its landed state, and a run linking to /runs/runone with its time 2026-10-24 11:00, read with no current change set; unknown IDs return 404 and POST 405. Read the diff: run.ls fields now carry json tags incl. started_at (quality f1), the test asserts the run time (f2). In a scratch copy, dropping the started_at tag makes the test fail with the run row's time empty, so the fix is covered.

## Links

### Out

- verifies [web-change-set-page](../capability/rhiza/web-monitor/web-change-set-page.md)

## History

- `01M6PFH625P3BYQAB8P6J38BBT` Web Roadmap pages: milestone groups and the backlog, and a change set's page with its intent, state, leaves, questions, runs and history
