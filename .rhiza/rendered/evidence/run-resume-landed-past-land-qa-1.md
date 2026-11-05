# run-resume-landed-past-land-qa-1

evidence · `01M6JWEB43RMMCMSTDJK3PHK4F` · actor [claude-code-3528205f](../actor/claude-code-3528205f.md)

## result

PASS

## statement

At 7c67207e I built rhiza with go build into a scratch dir and drove it by hand in scratch repos (RHIZA_* unset, fake rex on PATH). A shape workflow (land -> wait-main -> close-out) run stopped at wait-main; I set its journal workflow to ship-change, ran rhiza land (change set landed), dropped the stopped entry and appended a dead-pid process. The merge-base binary's rhiza run resume --foreground refused it (exit 3, 'is landed, no longer accepted', next rhiza run abandon). The built binary resumed it: wait-main continuing visit 1, pass, close-out ran (wrote its file), 'Run ... succeeded', exit 0. A separate ship-change run on a landed change set whose journal only entered build was still refused with exit 3 naming rhiza run abandon, journal unchanged at 3 lines.

## Links

### Out

- verifies [run-resume-landed-past-land](../capability/rhiza/workflows/workflow-runs/run-resume-landed-past-land.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
