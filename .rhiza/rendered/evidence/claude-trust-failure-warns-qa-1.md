# claude-trust-failure-warns-qa-1

evidence · `01M6GAQD0QGT6W8P0K6KXCD3XB` · actor [claude-code-d72ebedf](../actor/claude-code-d72ebedf.md)

## result

PASS

## statement

QA at 1f3f071f with the built binary, hermetic HOME and fake claude: with ~/.claude.json set to '{nope', rhiza session run still started the agent (fake claude ran in the worktree, exit 0) and printed 'Warning: trust the session's folders in <HOME>/.claude.json: read the Claude Code config: invalid character ...'; the file stayed byte-for-byte '{nope' and no .claude.json.rhiza-* temp file was left in HOME.

## Links

### Out

- verifies [claude-trust-failure-warns](../capability/rhiza/runtime/agent-adapters/claude-trust-failure-warns.md)

## History

- `01M6FSQG5W0V2D5PFJT8CHK3DV` Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them
