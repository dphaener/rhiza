# roadmap-across-projects-verified

evidence · `01M6ST3PPXKK6VTNRG8ZN341PX` · actor [claude-code-51ab7b84](../actor/claude-code-51ab7b84.md)

## result

PASS

## statement

At 318b2ef1, go test -trimpath ./cmd/rhiza -run 'TestScripts/(roadmap|describe-golden|inbox|health)' passes, including the new roadmap-projects.txtar: with a root of two projects (alpha with a milestone, beta with a backlog change set), rhiza roadmap --json at the root gives alpha's milestone group and beta's backlog group, each group and change set tagged with its project; --project alpha/beta and running from inside alpha give only that project's part, with no project tags; --project nope fails with 'not registered'. By hand with a scratch build and a fresh two-project root, the human roadmap shows '## alpha: ma' and '## beta: Backlog', and --limit 1 paging continues from alpha onto beta's backlog, ending with next_cursor null. go test ./internal/api -run 'Roadmap|Health|Milestone|Plan|Inbox' passes. The usage/limit tests in internal/api pass at this commit in a clean worktree; they failed only inside the verifier checkout, and also fail there on code this diff does not touch. Read the diff: health and the web milestones view keep a local-only roadmap; cursors carry the project.

## Links

### Out

- verifies [roadmap-across-projects](../capability/rhiza/orchestration/roadmap-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
