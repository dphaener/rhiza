# answer-6cd5e83f-3

decision · `01M6FWHFCNF02N323EY0DY5AAK`

## options

```json
[
  "Main root + session dir (Recommended)",
  "Main checkout root only"
]
```

## picked

Main root + session dir (Recommended)

## question

Which paths should be marked trusted in ~/.claude.json (or $CLAUDE_CONFIG_DIR/.claude.json)?

## session

6cd5e83f-56e1-4870-995e-01a19f9be0a3

## statement

Which paths should be marked trusted in ~/.claude.json (or $CLAUDE_CONFIG_DIR/.claude.json)? — Main root + session dir (Recommended)

## wording

the agent's option, picked by the owner in session 6cd5e83f-56e1-4870-995e-01a19f9be0a3

## Links

### Out

- derived-from [trust-fix-plan](../source/trust-fix-plan.md)

## History

- `01M6FSQG5W0V2D5PFJT8CHK3DV` Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them
