# custom-workflow-doc-reference-plain-path

decision · `01M6ZJ34KYNGHQ7ZWYA5A8X97G`

## reasoning

All three attempts failed only on the reference link. docs/reference/workflow-file.md does not exist: rhiza docs writes it in the docs-reference-workflow leaf of the accepted but unlanded Generated reference change set. A real link now fails internal/doclinks TestRepoDocsLinksResolve, so no builder can meet the literal statement, and builders 2 and 3 correctly refused to guess at the open question custom-workflow-doc-reference-link. Making this docs change set depend on Generated reference would hold the whole Agents, runs and daemon docs set back for one link, and the reference is already reachable through rhiza guide workflows, which the doc also names. The rest of the doc (the bugfix example, the check and describe walkthrough) is done. The producer of a file is the natural owner of links to it, so the plain-text path plus a follow-up link in Generated reference keeps doclinks green and the leaf's intent intact. This also settles question custom-workflow-doc-reference-link the same way (its second option).

## statement

docs-custom-workflow ships with the generated workflow-file reference named in plain text (docs/reference/workflow-file.md, alongside rhiza guide workflows), not as a Markdown link. That satisfies the leaf's 'links the generated workflow-file reference' for this change set, and verifiers accept it as is. The Generated reference change set (leaf docs-reference-workflow) turns the plain path into a link in docs/custom-workflows.md when it writes the file. This change set does not take a dependency on Generated reference.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [leaf-failed-01m6zh70-1](../question/leaf-failed-01m6zh70-1.md)
- answers [custom-workflow-doc-reference-link](../question/custom-workflow-doc-reference-link.md)

## History

- `01M6W5ETVXHS9EJQ2KH7STXY50` Agents, runs and daemon docs: docs/ explains using Rhiza with coding agents, workflows and runs, the daemon and autonomy, schedules and writing a custom workflow, with Rex and Claude Code setup
