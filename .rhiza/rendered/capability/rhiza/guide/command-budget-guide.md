# command-budget-guide

capability · `01M76XY4EE4QA5EKNQXHR7R9XH` · parent [guide](../guide.md)

## acceptance

go test -trimpath ./internal/guide -run 'Build|Review' asserts build.md and review.md each hold the budget rule (no sleep or poll, about four minutes per command, background once and wait once) and that it names no Go, make or Rhiza-repo command.

## statement

The build and review guide topics tell builders, verifiers, checks and fixers never to sleep or poll, never to run a single command longer than about four minutes, and to run a long build or test once in the background and wait for it once, in language-neutral words.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../decision/efficient-cycle-milestone-shape.md)

### In

- verifies [command-budget-guide-verified](../../../evidence/command-budget-guide-verified.md)
- verifies [command-budget-guide-qa-1](../../../evidence/command-budget-guide-qa-1.md)
- verifies [pr-501-command-budget](../../../evidence/pr-501-command-budget.md)

## History

- `01M76XQX53GGNWVCQHXVV510NH` Builder, verifier, check and fixer guidance says: never sleep or poll, no single command longer than about four minutes, and a long command runs once in the background and is waited for once (afk execute.js rule). Evidence: verifier session p90 33.7m against a 3.7m median, builder p90 25-36m against about 4m (prompt-gap rec 6).
