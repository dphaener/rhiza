# workflow-docs-examples-verified

evidence · `01M77YK07VPER9T9MTJAFQFQS5` · actor [claude-code-91d63dcc](../actor/claude-code-91d63dcc.md)

## result

PASS

## statement

At aa5d5f235, go test -trimpath -count=1 ./cmd/rhiza -run TestDocsExampleWorkflows passes: it globs docs/examples/workflows/*.toml (go, ruby, javascript, python; fails if fewer than 4) and runs rhiza workflow check on each, requiring exit 0 and 'is valid'. A scratch build's rhiza workflow check on each of the four files prints 'is valid', exit 0; mutated copies (a RHIZA_ env name, a dangling on.pass) fail with exit 1, so the check is real. docs/custom-workflows.md gains a 'Cleanup and environment' section explaining [[cleanup]] (runs once at run end in the run worktree, failure never changes the result) and [env] ({run_tmp}, RHIZA_* refused), and links the four examples, each with build, test, [env] and [[cleanup]].

## Links

### Out

- verifies [workflow-docs-examples](../capability/rhiza/workflows/workflow-authoring/workflow-docs-examples.md)

## History

- `01M763Z8TGV8GJK2E621JS6HYZ` Language-neutral workflow scaffold, docs and examples: rhiza workflow new writes a placeholder skeleton with extensive comments and example commands for popular languages; the docs explain the cleanup phase and workflow environment and carry example workflows for several languages.
