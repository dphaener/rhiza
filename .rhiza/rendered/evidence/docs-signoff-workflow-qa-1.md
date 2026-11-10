# docs-signoff-workflow-qa-1

evidence · `01M7398V2HMGJ7A150RSGTVYBY` · actor [claude-code-32c9dfde](../actor/claude-code-32c9dfde.md)

## result

PASS

## statement

QA at 0e35c5a8 with the branch binary in this repo: rhiza workflow check docs is valid; workflow describe docs shows 'autonomy signoff owner' (and --json autonomy [{stage signoff, mode owner}]), includes docs-build and docs-review; docs-build's builder and docs-review's fixer prompts name the humanize skill and the doc-prose-humanized / doc-examples-plain-commands constraints; change show --json of README and quickstart polish, Concepts polish and Agents and runs docs polish each shows workflow docs; make self-check passes in a clean worktree at HEAD.

## Links

### Out

- verifies [docs-signoff-workflow](../capability/rhiza/documentation/docs-signoff-workflow.md)

## History

- `01M71G46Q5378CYYVEHFF2EMP3` Docs sign-off workflow: a repo workflow named docs (ship-change with signoff = owner, builders and fixers told to write with the humanize skill) builds doc change sets, and rhiza change new --workflow and a command for an existing draft pin a change set's workflow
