# workflow-new-placeholder

capability · `01M7642ND4H8MY27X1Y6D2CZ2K` · parent [workflow-authoring](../workflow-authoring.md)

## acceptance

The workflow-new testscript shows workflow new writing the skeleton (exit 0) with no make test and the examples present, workflow check on it exiting 0 with a warning naming the placeholder, and no such warning once the placeholder is replaced.

## statement

rhiza workflow new without --from writes a skeleton whose test command is a placeholder, with comments on every field and example commands for Go, Ruby, JavaScript/TypeScript, Python and Rust, plus a commented cleanup step and environment example; the placeholder is a workflow check warning, which never makes the workflow invalid, so workflow new still validates and writes the skeleton, and workflow check warns about it until it is replaced.

## Links

### Out

- decided-by [d-scaffold-placeholder](../../../../decision/d-scaffold-placeholder.md)

### In

- blocks [review-01m75w8yt6m1f8w05gajzsec0h-r3-f2](../../../../question/review-01m75w8yt6m1f8w05gajzsec0h-r3-f2.md)
- verifies [workflow-new-placeholder-verified](../../../../evidence/workflow-new-placeholder-verified.md)
- verifies [workflow-new-placeholder-qa-1](../../../../evidence/workflow-new-placeholder-qa-1.md)
- verifies [pr-490-language-neutral-workflow](../../../../evidence/pr-490-language-neutral-workflow.md)

## History

- `01M763Z8TGV8GJK2E621JS6HYZ` Language-neutral workflow scaffold, docs and examples: rhiza workflow new writes a placeholder skeleton with extensive comments and example commands for popular languages; the docs explain the cleanup phase and workflow environment and carry example workflows for several languages.
