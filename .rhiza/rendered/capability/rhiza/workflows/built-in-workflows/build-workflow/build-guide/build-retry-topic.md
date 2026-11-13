# build-retry-topic

capability · `01M7974G5P60MM2EHVR5M79X51` · parent [build-guide](../build-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a build-retry test that finds the inspect-the-working-tree rule, rhiza run status --json, address every finding as written, the same `Leaf:` trailer fix commit and never rewrite history, and asserts the text stays under a third of the TopicCap and repeats none of the builder topic's shared blocks.

## statement

`rhiza guide build-retry` tells a builder on a later attempt only what a retry adds: inspect the working tree and the earlier attempts' commits for their partial work, read the previous attempt's findings (in the prompt when there, else `rhiza run status --json`), address every finding exactly as written, and add one fix commit with the same Leaf: trailer, never rewriting history.

## Links

### Out

- depends-on [guide-shared-blocks](guide-shared-blocks.md)
- decided-by [build-role-topics-shape](../../../../../../decision/build-role-topics-shape.md)

### In

- depends-on [build-role-skills-installed](build-role-skills-installed.md)
- verifies [build-retry-topic-verified](../../../../../../evidence/build-retry-topic-verified.md)
- verifies [build-retry-topic-qa-1](../../../../../../evidence/build-retry-topic-qa-1.md)
- verifies [build-retry-topic-qa-2](../../../../../../evidence/build-retry-topic-qa-2.md)
- verifies [pr-521-build-role-topics](../../../../../../evidence/pr-521-build-role-topics.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
