# resume-landed-rebase-resolution

decision · `01M6K3KEAXMJ33NJB0XGREM50V`

## statement

Resolve the internal/api/daemon.go conflict in daemonPoll1 by keeping this branch's daemon.Resumable call (change-set state lookup and planWorkflow, for the landed-past-land exception) under main's renamed guard `!stopped && !paused` (main renamed low to stopped when the usage governor landed); merge origin/main into the branch rather than rebase and re-render

## wording

the orchestrator's, deciding as the owner's proxy under the owner's standing instruction "act on my behalf" (2026-11-03); assumed, not the owner's words

## Links

### Out

- answers [rebase-conflict-fix-a-run-interrupted-after](../question/rebase-conflict-fix-a-run-interrupted-after.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
