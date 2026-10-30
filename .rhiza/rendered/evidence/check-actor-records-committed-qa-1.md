# check-actor-records-committed-qa-1

evidence · `01M653KST9SDXW4DF6M48J2M87` · actor [claude-code-2f55ab6c](../actor/claude-code-2f55ab6c.md)

## result

PASS

## statement

Built the branch head (go build ./cmd/rhiza into the scratchpad) and ran a real `rhiza run start --workflow review` in a fresh scratch repo (RHIZA_* unset, own HOME), with only rex and claude faked: check alpha recorded an empty findings round and wrote no model entry; check beta recorded an empty round and added Evidence. The run succeeded. Alpha's idle session got commit 'chore(rhiza): record alpha check writes' with trailers 'Change: <cs>' and 'Actor: claude-code-84fb297e', no Leaf:, touching only its actor's 3 record entries under .rhiza/log/; beta's commit held its Evidence plus its record entries. Afterwards `git status --porcelain --untracked-files=all -- .rhiza` in the run worktree was empty. The no-uncommitted-means-no-commit path was checked via go test ./internal/api -run TestCommitCheckWrites (both PASS); go test ./cmd/rhiza -run TestScripts/review- ok.

## Links

### Out

- verifies [check-actor-records-committed](../capability/check-actor-records-committed.md)

## History

- `01M602TH4C724WZB9MKBN5Z1PP` A review check or critic session that writes no model entries still gets its actor's record entries committed, so no uncommitted record change set is left in the run worktree for a later fixer to sweep up
