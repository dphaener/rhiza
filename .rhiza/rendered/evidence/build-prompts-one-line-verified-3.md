# build-prompts-one-line-verified-3

evidence · `01M79BXW5TVDTSJQXW80BDWSYY` · actor [claude-code-de7c7c9f](../actor/claude-code-de7c7c9f.md)

## result

PASS

## statement

Attempt 3 verifier at e78b3b302 (detached). Re-ran every Checks: command of e78b3b302 and 41132e7d0: go test -trimpath ./internal/workflow -run 'TestBuiltinRolePromptsScopeTests|TestBuildBuiltin' (exit 0), go test -trimpath ./internal/guide -run 'TestBuilderText|TestVerifierText' (exit 0), rhiza ship generate and rhiza render with this commit's build (exit 0, no tracked diff), make self-check (exit 0, only pre-existing unclaimed warnings). Re-ran the leaf acceptance: go test -trimpath ./internal/workflow/... ok (incl. the docs-build/docs-review test at workflow_test.go:2420), go test -trimpath ./cmd/rhiza -run 'TestScripts/(build-evidence-commit|docs-commands-exist|describe-golden)$' ok. Read build.toml:18/24, review.toml:67, .rhiza/workflows/docs-build.toml:21/27 and docs-review.toml:70: the one-line builder/verifier prompts match the statement byte for byte, docs-build keeps its humanize addendum. Model findings f1/f2 are addressed: role-prompts-scoped-tests now exempts the one-line /rhiza-builder and /rhiza-verifier prompts (rule comes from the topics' shared test-scope block), and builder-prehandoff-checklist names the builder topic and TestBuilderText and drops the prompt-points-at-checklist clause.

## Links

### Out

- verifies [build-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-prompts-one-line.md)

## History

- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
