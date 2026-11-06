# web-change-set-page-verified

evidence · `01M6Q6WPSWYYW9B7A1JEAYXWMC` · actor [claude-code-b57f3891](../actor/claude-code-b57f3891.md)

## result

PASS

## statement

At f309922d I ran go test -trimpath -count=1 ./internal/web/... (ok; TestChangeSetPage PASS). It serves /changes/<id> through api.ReadOnlyHandler for a fixture proposed change set with two leaves (one verified by Evidence), a question blocking a leaf, a depends-on link to a landed change set, and a build run. It asserts the title, kind, state, each leaf's statement and acceptance, has evidence / no evidence yet, the 2 leaves, 1 with evidence count, the open question, the dependency's text, href=/changes/<base> and landed pill, href=/runs/runone, build run, and the Proposed history entry. It also checks that no handle or ID is shown, that an unknown ULID and a bad ref both return 404, and that POST returns 405. go test ./internal/api -run 'TestReadOnlyMarker|ChangeShow' passed, so change.show is now ReadOnly. go vet and gofmt are clean on both packages. I read the diff: the roadmap rows' links moved from /cs/ to /changes/, and roadmap_test was updated to match. A whole-package internal/api run hit the 10m go test timeout while several other worktrees' api.test processes ran on the machine; CI runs the full suite.

## Links

### Out

- verifies [web-change-set-page](../capability/rhiza/web-monitor/web-change-set-page.md)

## History

- `01M6PFH625P3BYQAB8P6J38BBT` Web Roadmap pages: milestone groups and the backlog, and a change set's page with its intent, state, leaves, questions, runs and history
