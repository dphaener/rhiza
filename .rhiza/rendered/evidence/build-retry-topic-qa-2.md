# build-retry-topic-qa-2

evidence · `01M79C8N24RXCY21Y2VPPVZ9ZR` · actor [claude-code-9b656f6b](../actor/claude-code-9b656f6b.md)

## result

PASS

## statement

`rhiza guide build-retry` (built at HEAD) prints a 192-token topic: inspect the working tree and earlier commits, read findings from the prompt or rhiza run status --json, address every finding exactly as written, one fix commit with the same Leaf: trailer, never rewrite history. `rhiza guide ls` lists it. go test -trimpath ./internal/guide/... ok

## Links

### Out

- verifies [build-retry-topic](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-retry-topic.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
