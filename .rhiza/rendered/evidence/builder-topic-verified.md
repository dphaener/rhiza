# builder-topic-verified

evidence · `01M799WYJXR6C9WE81WEWAW5JS` · actor [claude-code-becb034d](../actor/claude-code-becb034d.md)

## result

PASS

## statement

Verified builder-topic at c04eaaa75 (attempt 2) in a detached checkout. Re-ran the commit's Checks: go test -trimpath -count=1 ./internal/guide/... (ok), ./internal/api -run Guide (ok), ./cmd/rhiza -run TestScripts/guide (ok), rhiza ship generate (exit 0) and rhiza render (exit 0, no tree changes). A scratch build's `rhiza guide builder | grep '^#'` now prints exactly the six ## sections in order (Task, Expected outcome, Read first, Must do, Must not do, Context): command-budget and outside-work sit as ### under Must do, never-list as ### under Must not do. That fixes attempt 1's finding. Read the diff: Expand's new `nested` marker demotes block headings and only builder.md uses it; TestSharedNestedBlockDemotesHeadings covers it. TestBuilderText asserts the only ## headings are the six, the shared session-done/test-scope/unstable-ci/command-budget/outside-work/never-list blocks expanded and placed in their sections, the read-first and question-before-work rules, acceptance commands, rhiza ship generate, rhiza render, a Checks: block above the single final trailer paragraph, the actor-record git add line, and the RHIZA_ATTEMPT > 1 rule naming /rhiza-build-retry and `rhiza guide build-retry`. A negative case bans go test/go build/go vet/make/testscript.

## Links

### Out

- verifies [builder-topic](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/builder-topic.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
