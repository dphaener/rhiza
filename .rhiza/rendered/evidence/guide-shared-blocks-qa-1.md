# guide-shared-blocks-qa-1

evidence · `01M79B54NT5M0ZKSB9ZTB0FHPE` · actor [claude-code-fb58129d](../actor/claude-code-fb58129d.md)

## result

PASS

## statement

Built rhiza from the branch head; ran `rhiza guide <topic>` for every topic in `rhiza guide ls`: no `<!-- shared:` marker survives and the builder/verifier/build-retry texts show the shared done, test-scope, unstable-CI, command-budget, outside-work and never-list blocks expanded in place; `rhiza guide nosuch` exits 4. go test -trimpath ./internal/guide/... ok.

## Links

### Out

- verifies [guide-shared-blocks](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/guide-shared-blocks.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
