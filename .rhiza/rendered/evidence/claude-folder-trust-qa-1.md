# claude-folder-trust-qa-1

evidence · `01M6GAQ795W66JYYEZXVV8VPXP` · actor [claude-code-d72ebedf](../actor/claude-code-d72ebedf.md)

## result

PASS

## statement

QA at 1f3f071f with the built binary (go build ./cmd/rhiza) in a scratch git repo, hermetic HOME, fake claude (dumps .claude.json at start) and fake rex: rhiza session start --permission-mode manual, then session run --dry-run left ~/.claude.json (numStartups, theme, projects./other.allowedTools) unchanged; the real session run set hasTrustDialogAccepted=true for the main checkout root (repo) and the session worktree (repo.qa-trust) before the agent started (the fake claude saw both), kept numStartups, theme and /other.allowedTools and the file mode 0640, left no temp file. With CLAUDE_CONFIG_DIR set to an empty dir, session run created $CLAUDE_CONFIG_DIR/.claude.json (0600) with both paths trusted and left HOME's file byte-identical. go test -run 'Trust|CommandWritesNoFile' for internal/api and claudecode passes (resume re-trust and Command writing no file are covered there). The owner's real ~/.claude.json uses the same hasTrustDialogAccepted key under projects.

## Links

### Out

- verifies [claude-folder-trust](../capability/rhiza/runtime/agent-adapters/claude-folder-trust.md)

## History

- `01M6FSQG5W0V2D5PFJT8CHK3DV` Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them
