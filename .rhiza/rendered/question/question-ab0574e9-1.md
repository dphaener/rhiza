# question-ab0574e9-1

question · `01M796WYTNKDNYWDP070HT3QQF`

## session

ab0574e9-12c9-4b2a-af3f-01a1d26b9855

## statement

Build role topics: how are the shared blocks (done rule, test scope, unstable CI, command budget, outside work, never-list) written once across the builder, verifier and builder-retry topics (and later the review and plan topics)? / Build role topics: how do the builder, verifier and builder-retry topics reach run sessions, and what happens to today's `build` topic and rhiza-build skill? / Build role topics: how does a builder retry reach the build-retry topic? / Build role topics: Rhiza's own .rhiza/workflows/docs-build.toml copies the builder and verifier prompts with a doc-prose addendum. What happens to it?

## Links

### In

- answers [build-role-topics-shape](../decision/build-role-topics-shape.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
