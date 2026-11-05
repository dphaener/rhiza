# daemon-resumes-landed-past-land-qa-1

evidence · `01M6JWEKHE3F0YEEFZD0DJ7BFS` · actor [claude-code-3528205f](../actor/claude-code-3528205f.md)

## result

PASS

## statement

At 7c67207e I built rhiza with go build into a scratch dir and drove it by hand in scratch repos (RHIZA_* unset, fake rex on PATH). On the same interrupted ship-change run past land on a landed change set, rhiza daemon --permission-mode manual --once printed 'Resumed run <id>', journaled one daemon-resumed entry and spawned a real detached rhiza run resume --foreground, which went on past wait-main through close-out to a finished/succeeded entry with no 'no longer accepted' refusal. A second --once --json gave resumed: []. In another repo, --once resumed none of: a landed run that never entered land, a landed run already daemon-resumed and interrupted again, and a failed run after land.

## Links

### Out

- verifies [daemon-resumes-landed-past-land](../capability/rhiza/orchestration/daemon/daemon-resumes-landed-past-land.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
