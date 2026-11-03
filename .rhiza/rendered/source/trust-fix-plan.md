# trust-fix-plan

source · `01M6FWDTQ4VASTE7MD19JV886K`

## quotes

```json
[]
```

## raw_sha256

7d70d7aed3ae9ba11c83a364e4096cd019848d909ecafc5c640180ab78c97995

## session

6cd5e83f-56e1-4870-995e-01a19f9be0a3

## summary

Planned the Fix for Claude Code sessions stalling at the folder-trust prompt in a new repo. The owner picked the agent's recommended options throughout. The Fix serves m24-on-the-vm. The trust write goes in a new optional adapter interface that claudecode implements. session run calls it before every agent start and skips it on --dry-run, so Adapter.Command stays pure. It marks the main checkout root and the session's working dir as trusted (hasTrustDialogAccepted=true) in ~/.claude.json, or in $CLAUDE_CONFIG_DIR/.claude.json when that is set. The write keeps every other key and is atomic (temp file and rename). It never overwrites a malformed file. When the file can't be read or written, the session starts anyway and the failure is reported as a warning.

## title

Planning the folder-trust fix, with the owner

## Links

### In

- derived-from [claude-folder-trust](../capability/rhiza/runtime/agent-adapters/claude-folder-trust.md)
- derived-from [claude-trust-failure-warns](../capability/rhiza/runtime/agent-adapters/claude-trust-failure-warns.md)
- derived-from [answer-6cd5e83f-1](../decision/answer-6cd5e83f-1.md)
- derived-from [answer-6cd5e83f-2](../decision/answer-6cd5e83f-2.md)
- derived-from [answer-6cd5e83f-3](../decision/answer-6cd5e83f-3.md)
- derived-from [answer-6cd5e83f-4](../decision/answer-6cd5e83f-4.md)

## History

- `01M6FSQG5W0V2D5PFJT8CHK3DV` Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them
