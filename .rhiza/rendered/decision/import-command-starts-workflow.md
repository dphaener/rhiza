# import-command-starts-workflow

decision · `01M6J31RRCX75M5596X156W6BR`

## statement

rhiza import runs rhiza init when needed and starts a built-in import workflow whose agent steps do the judgment; the command is only the entry point.

## wording

Command starts an import workflow: rhiza import runs init if needed and starts a built-in import workflow: survey, record areas (foreach), decisions, Sources and retire docs, planned work as drafts, verify (rhiza check plus the repo's own tests). Judgment stays with agents; the command is the entry point. The assistant's option, picked by the owner.

## Links

### Out

- derived-from [milestone-26-interview](../source/milestone-26-interview.md)

### In

- decided-by [m26-import-repeatable](../intent/milestone-26/m26-import-repeatable.md)

## History

- `01M6J2RTBPHTT1VFJYW2ZA5FRC` Milestone 26: an existing repo with built code comes into Rhiza, Nebari first, and its work then runs on Rhiza workflows
