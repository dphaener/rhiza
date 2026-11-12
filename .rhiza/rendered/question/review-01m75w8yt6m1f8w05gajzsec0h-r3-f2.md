# review-01m75w8yt6m1f8w05gajzsec0h-r3-f2

question · `01M7656JXFA5K23CPA9C6E2TE0`

## statement

Plan review round 3 left finding f2 (contradiction on workflow-new-placeholder): Change set 01M763Z8TGV8GJK2E621JS6HYZ ('language-neutral workflow scaffold') writes 'a commented cleanup step and environment example' into the scaffold, and workflow-docs-examples documents 'workflow environment'. Workflow env only exists through workflow-env in change set 01M763YX47S5XY0W3P8FVJ54S8. Z8TG depends only on the cleanup-phase change set 01M763YS08ZM13JZ065Q95YA2E, so it can be built before workflow env exists and would scaffold and document a field that isn't there yet. It needs a depends-on link to 01M763YX47S5XY0W3P8FVJ54S8. The critic confirmed it: In the log, 01M763Z8TGV8GJK2E621JS6HYZ has only one depends-on link, to 01M763YS08ZM13JZ065Q95YA2E. Yet workflow-new-placeholder scaffolds an environment example and workflow-docs-examples documents workflow environment. Workflow env only comes from workflow-env in 01M763YX47S5XY0W3P8FVJ54S8, so Z8TG really does need YX47's work. How should it be settled?

## Links

### Out

- blocks [workflow-new-placeholder](../capability/rhiza/workflows/workflow-authoring/workflow-new-placeholder.md)

### In

- answers [scaffold-env-example-waits](../decision/scaffold-env-example-waits.md)
- answers [answer-review-01m75w8yt6m1f8w05gajzsec0h-r3-f2](../decision/answer-review-01m75w8yt6m1f8w05gajzsec0h-r3-f2.md)

## History

- `01M763Z8TGV8GJK2E621JS6HYZ` Language-neutral workflow scaffold, docs and examples: rhiza workflow new writes a placeholder skeleton with extensive comments and example commands for popular languages; the docs explain the cleanup phase and workflow environment and carry example workflows for several languages.
