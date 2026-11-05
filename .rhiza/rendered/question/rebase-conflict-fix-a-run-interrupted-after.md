# rebase-conflict-fix-a-run-interrupted-after

question · `01M6K3ACXBM7B6F5Q0EJ9HZP17`

## statement

Rebasing branch cs/fix-a-run-interrupted-after onto origin/main conflicts outside the rendered files, so ship-change stopped. How should these conflicts be resolved: internal/api/daemon.go?

## Links

### Out

- blocks change set `01M6J965BRJ2ZTDT7Y311PM96G` (Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless)

### In

- answers [resume-landed-rebase-resolution](../decision/resume-landed-rebase-resolution.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
