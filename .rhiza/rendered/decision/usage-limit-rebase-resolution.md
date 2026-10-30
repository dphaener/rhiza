# usage-limit-rebase-resolution

decision · `01M65PE1MKZ9N5X8M52H9ZDR0N`

## statement

The usage-limit branch conflicts with main only where both sides append to the same test files and to rendered index; both sides are kept, the branch's UsageLimit stubs after main's methods, and the rendered files are re-rendered

## wording

Arbiter answer under the standing autonomous instruction: keep both sides, re-render; go vet and the internal/api tests pass

## Links

### Out

- answers [rebase-conflict-usage-limit-pause](../question/rebase-conflict-usage-limit-pause.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
