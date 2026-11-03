# pr-205-fix

evidence · `01M6GBDHHVK439XB3T18V8ARRJ` · actor [dphaener](../actor/dphaener.md)

## head

7eef89cee973c315264eb04c198a1041e270554b

## pr

https://github.com/dphaener/rhiza/pull/205

## result

1 check passed: check SUCCESS

## statement

Pull request #205 (https://github.com/dphaener/rhiza/pull/205) passed CI at 7eef89c

## status

verified

## Links

### Out

- verifies [claude-folder-trust](../capability/rhiza/runtime/agent-adapters/claude-folder-trust.md)
- verifies [claude-trust-failure-warns](../capability/rhiza/runtime/agent-adapters/claude-trust-failure-warns.md)

## History

- `01M6FSQG5W0V2D5PFJT8CHK3DV` Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them
