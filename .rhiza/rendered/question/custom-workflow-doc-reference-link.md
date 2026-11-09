# custom-workflow-doc-reference-link

question · `01M6ZHSA0Q992EDPRM65Q12PEP`

## statement

@docs-custom-workflow must link the generated workflow-file reference, docs/reference/workflow-file.md, but that file comes from the docs-reference-workflow leaf of the Generated reference change set (01M6W5E3JQXD62WN4RHG546KRX), which is accepted and not landed, and this change set does not depend on it. A link now fails internal/doclinks TestRepoDocsLinksResolve. Should this change set depend on Generated reference and add the link once it lands, or ship with the path named in plain text and the link added by Generated reference when it writes the file?

## Links

### Out

- blocks [docs-custom-workflow](../capability/rhiza/documentation/docs-custom-workflow.md)

### In

- answers [custom-workflow-doc-reference-plain-path](../decision/custom-workflow-doc-reference-plain-path.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
