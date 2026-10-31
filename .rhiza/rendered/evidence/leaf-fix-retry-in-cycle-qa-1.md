# leaf-fix-retry-in-cycle-qa-1

evidence · `01M688W60QPHE0XKJ213TMN9PK` · actor [claude-code-96da0913](../actor/claude-code-96da0913.md)

## result

PASS

## statement

Built rhiza at 04cd35f (go build -o $TMPDIR/rhiza ./cmd/rhiza) and ran the built-in review workflow with rhiza run start --workflow review --permission-mode manual in hermetic scratch repos (own HOME and git config; fake rex running the pane command, fake claude acting each role) on a built two-leaf change set where quality's cycle-1 major on leaf alpha is confirmed. Fixer exits 1 on its first session: cycle 2 ran fixer, fixer, verifier, then quality check 2; journal leaf-attempt fixer 2 fail, fixer 3 pass, verifier 3 pass, all before the cycle-2 check session; quality passed after 2 cycles and the run reached sign-off. Verifier exits 1 on its first session: cycle 2 ran fixer, verifier, fixer, verifier (a fresh pair), then quality check 2; passed after 2 cycles, run at sign-off.

## Links

### Out

- verifies [leaf-fix-retry-in-cycle](../capability/fixer-retry-in-cycle/leaf-fix-retry-in-cycle.md)

## History

- `01M6580HBAKHN7XTR7XQV3XR0Q` A failed review leaf fixer, or its verifier, is retried within its cycle as the change-set fixer is, so only a leaf fix failing every attempt spends the cycle
