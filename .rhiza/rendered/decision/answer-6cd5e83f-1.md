# answer-6cd5e83f-1

decision · `01M6FWGTK66ESQJCG9Q4310ESN`

## options

```json
[
  "m24-on-the-vm (Recommended)",
  "m24-headless-runs"
]
```

## picked

m24-on-the-vm (Recommended)

## question

Which intent should the trust-prompt Fix serve? It serves none yet. m24-on-the-vm (Rhiza runs on the new VM, where each repo is new to Claude) fits best.

## session

6cd5e83f-56e1-4870-995e-01a19f9be0a3

## statement

Which intent should the trust-prompt Fix serve? It serves none yet. m24-on-the-vm (Rhiza runs on the new VM, where each repo is new to Claude) fits best. — m24-on-the-vm (Recommended)

## wording

the agent's option, picked by the owner in session 6cd5e83f-56e1-4870-995e-01a19f9be0a3

## Links

### Out

- derived-from [trust-fix-plan](../source/trust-fix-plan.md)

## History

- `01M6FSQG5W0V2D5PFJT8CHK3DV` Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them
