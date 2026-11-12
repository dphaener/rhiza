# command-budget-guide-qa-1

evidence · `01M789356H4YEZJK26J6GY7HM0` · actor [claude-code-9c9e773f](../actor/claude-code-9c9e773f.md)

## result

PASS

## statement

QA at e629ee6d0: built the binary (go build -trimpath ./cmd/rhiza) and ran `rhiza guide build` and `rhiza guide review`: each prints '## Command budget' with 'Never sleep and never poll. No command runs longer than about four minutes; a longer build or test runs once in the background and is waited for once.', naming no language, make or Rhiza command. Ran `rhiza guide install --agent claude` into a temp HOME/repo: the installed rhiza-build and rhiza-review SKILL.md both carry the section. go test -trimpath -count=1 ./internal/guide passes (budget test and token caps).

## Links

### Out

- verifies [command-budget-guide](../capability/rhiza/guide/command-budget-guide.md)

## History

- `01M76XQX53GGNWVCQHXVV510NH` Builder, verifier, check and fixer guidance says: never sleep or poll, no single command longer than about four minutes, and a long command runs once in the background and is waited for once (afk execute.js rule). Evidence: verifier session p90 33.7m against a 3.7m median, builder p90 25-36m against about 4m (prompt-gap rec 6).
