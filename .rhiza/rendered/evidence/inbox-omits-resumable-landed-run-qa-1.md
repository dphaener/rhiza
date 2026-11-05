# inbox-omits-resumable-landed-run-qa-1

evidence · `01M6JWETQQFAMRCJQJ1VKNG3VF` · actor [claude-code-3528205f](../actor/claude-code-3528205f.md)

## result

PASS

## statement

At 7c67207e I built rhiza with go build into a scratch dir and drove it by hand in scratch repos (RHIZA_* unset, fake rex on PATH). rhiza inbox on the interrupted ship-change run past land on a landed change set printed 'Nothing waits on you' (the merge-base binary listed it as skipped with rhiza run abandon). In a repo with three landed change sets and fixture journals, rhiza inbox listed exactly: skipped run N-LAND (never reached land) with next rhiza run abandon N-LAND, failed run F-LAND (main-green failed after land) with next rhiza run resume --retry, and gave-up run G-LAND (entered land, daemon-resumed, interrupted again) with next rhiza run resume G-LAND.

## Links

### Out

- verifies [inbox-omits-resumable-landed-run](../capability/rhiza/orchestration/inbox/inbox-omits-resumable-landed-run.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
