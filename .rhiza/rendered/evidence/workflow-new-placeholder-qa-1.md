# workflow-new-placeholder-qa-1

evidence · `01M77ZTTG43QVX4Q7X43W9Z5M8` · actor [claude-code-65a2e8fe](../actor/claude-code-65a2e8fe.md)

## result

PASS

## statement

Built rhiza at the branch head and, in a scratch git repo after rhiza init: workflow new demo exited 0 and wrote a skeleton with no make test, comments on every item, example commands for Go, Ruby, JavaScript/TypeScript, Python and Rust, and commented [[cleanup]] and [env] examples; the placeholder command run by sh prints its message and exits 1; workflow check demo exited 0 with 'warning: step test still runs the placeholder command from workflow new' (--json valid true with the warning); uncommenting the cleanup and env examples still checks valid; after sed replaced run with cargo test, workflow check exited 0 with no warning; a --from ship-change fork has no placeholder warning. go test ./cmd/rhiza -run TestScripts/workflow-new passes.

## Links

### Out

- verifies [workflow-new-placeholder](../capability/rhiza/workflows/workflow-authoring/workflow-new-placeholder.md)

## History

- `01M763Z8TGV8GJK2E621JS6HYZ` Language-neutral workflow scaffold, docs and examples: rhiza workflow new writes a placeholder skeleton with extensive comments and example commands for popular languages; the docs explain the cleanup phase and workflow environment and carry example workflows for several languages.
