# workflow-docs-examples-qa-1

evidence · `01M77ZTTP5AJMZ3YRN32WK1NFR` · actor [claude-code-65a2e8fe](../actor/claude-code-65a2e8fe.md)

## result

PASS

## statement

Read docs/custom-workflows.md: it has a Cleanup and environment section on [[cleanup]] and [env] ({run_tmp}, RHIZA_* refused) and links Go, Ruby, JavaScript/TypeScript and Python examples. Ran the built binary's workflow check on each of docs/examples/workflows/{go,ruby,javascript,python}.toml: each exited 0, 'is valid', no warnings. go test ./cmd/rhiza -run TestDocsExampleWorkflows passes.

## Links

### Out

- verifies [workflow-docs-examples](../capability/rhiza/workflows/workflow-authoring/workflow-docs-examples.md)

## History

- `01M763Z8TGV8GJK2E621JS6HYZ` Language-neutral workflow scaffold, docs and examples: rhiza workflow new writes a placeholder skeleton with extensive comments and example commands for popular languages; the docs explain the cleanup phase and workflow environment and carry example workflows for several languages.
