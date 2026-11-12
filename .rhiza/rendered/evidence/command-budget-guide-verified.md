# command-budget-guide-verified

evidence · `01M788XRBD3HZQGS3G7M7WTDMP` · actor [claude-code-ee453ae7](../actor/claude-code-ee453ae7.md)

## result

PASS

## statement

At ec0f384fe ran go test -trimpath -count=1 ./internal/guide -run 'Build|Review' -v: TestBuildReviewCommandBudget PASS, asserting build.md and review.md each hold '## Command budget' with 'Never sleep and never poll', 'about four minutes' and 'runs once in the background and is waited for once', and that the section names no go/make/rhiza/testscript/Claude. Read the diff: the section is language-neutral in both files. The attempt-1 finding (TestCaps over cap) is fixed: full ./internal/guide passes, build.md 3175 bytes/794 tokens, review.md 3200 bytes/800 tokens, both within the 800 cap; the trims keep every phrase the other guide tests assert. Scoped testscripts guide, guide-install, build-evidence-commit and docs-commands-exist pass.

## Links

### Out

- verifies [command-budget-guide](../capability/rhiza/guide/command-budget-guide.md)

## History

- `01M76XQX53GGNWVCQHXVV510NH` Builder, verifier, check and fixer guidance says: never sleep or poll, no single command longer than about four minutes, and a long command runs once in the background and is waited for once (afk execute.js rule). Evidence: verifier session p90 33.7m against a 3.7m median, builder p90 25-36m against about 4m (prompt-gap rec 6).
