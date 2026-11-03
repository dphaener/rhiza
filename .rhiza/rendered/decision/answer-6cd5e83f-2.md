# answer-6cd5e83f-2

decision · `01M6FWH59C2EPHHMB0KC18AD4X`

## options

```json
[
  "Optional adapter method (Recommended)",
  "Inside Command"
]
```

## picked

Optional adapter method (Recommended)

## question

Where should the trust write go? Today Adapter.Command is pure (it only builds argv), and session run calls ctx.Start right after it.

## session

6cd5e83f-56e1-4870-995e-01a19f9be0a3

## statement

Where should the trust write go? Today Adapter.Command is pure (it only builds argv), and session run calls ctx.Start right after it. — Optional adapter method (Recommended)

## wording

the agent's option, picked by the owner in session 6cd5e83f-56e1-4870-995e-01a19f9be0a3

## Links

### Out

- derived-from [trust-fix-plan](../source/trust-fix-plan.md)

## History

- `01M6FSQG5W0V2D5PFJT8CHK3DV` Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them
