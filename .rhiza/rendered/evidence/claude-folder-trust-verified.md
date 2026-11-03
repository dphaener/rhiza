# claude-folder-trust-verified

evidence · `01M6G7AD4BCG45GR4558FW5C4V` · actor [claude-code-ce2adfaa](../actor/claude-code-ce2adfaa.md)

## result

PASS

## statement

Verifier at 4bd51c50: read the diff (adapter.Truster; claudecode TrustFile/Trust; session run calls trustFolders before every ctx.Start inside the pause loop, after the --dry-run early return; atomic temp-file rename keeping the file mode). Ran go test -count=1 -run 'Trust|CommandWritesNoFile' -v ./internal/api ./internal/adapter/claudecode: all PASS. With a hermetic temp HOME and fake adapters these show: TestSessionRunTrustsFolders, the main checkout root and the worktree set to hasTrustDialogAccepted=true, while numStartups, theme, /other's allowedTools and the worktree's lastCost stay unchanged and mode stays 0600; TestSessionRunTrustHonoursConfigDirAndCreatesFile, CLAUDE_CONFIG_DIR/.claude.json is created when missing and HOME's file is left alone; TestSessionRunResumeTrustsAgain, after the config is removed during a usage-limit pause, the resumed start trusts both paths again; TestSessionRunTrustDryRunWritesNothing, HOME stays empty; TestCommandWritesNoFile, Adapter.Command writes nothing; plus TestTrust/TestTrustFile (idempotent, big numbers kept, broken JSON is a warning and is not rewritten). make fmt-check, vet and self-check pass. In go test ./..., only the wall-clock read-budget tests failed (api TestReadBudget, slice TestBudget, modelapp TestGeneratedModel): this 4-core host ran at load average 10-15 with other worktrees' checks running, and the diff touches no read path.

## Links

### Out

- verifies [claude-folder-trust](../capability/rhiza/runtime/agent-adapters/claude-folder-trust.md)

## History

- `01M6FSQG5W0V2D5PFJT8CHK3DV` Fix: Claude Code sessions stop at the folder-trust prompt in a new git repo: before starting an interactive session, the claudecode adapter should set projects[<main checkout root>].hasTrustDialogAccepted=true in ~/.claude.json (as afk internal/provider/trust.go TrustFolder does), so new projects never need the owner to trust them
