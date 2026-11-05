# inbox-omits-resumable-landed-run-verified

evidence · `01M6JW3255MN1RAZFMCZKXY0M4` · actor [claude-code-86d99ea0](../actor/claude-code-86d99ea0.md)

## result

PASS

## statement

At 6a768be8, cmd/rhiza/testdata/script/inbox-landed-resumable.txtar runs rhiza inbox --json over fixture journals on four landed change sets and asserts all three acceptance cases: (1) R-LAND, interrupted after entering land, is absent; (2) G-LAND, resumed by the daemon then interrupted again, is gave-up, and F-LAND, failed at main-green after land, is failed; (3) N-LAND, an interrupted build run that never reached land, is skipped with '"next": "rhiza run abandon N-LAND"' in JSON and 'next: rhiza run abandon N-LAND' in text (the attempt-1 finding), with exactly 3 items. go test ./cmd/rhiza -run TestScripts/inbox -count=1 passes. Mutation: with internal/api/inbox.go restored to 3b04cbed, the script fails at line 19 (R-LAND listed). The diff in inboxItems makes gone() false for a landed change set when a non-plan run has Visits[land]>0, used for the failed, skipped and gave-up branches. fmt-check, vet and self-check pass; go test ./... passes except timing budgets (TestReadBudget, slice TestBudget) and a plan-workflow timeout under load 12-17 on 4 cores, which pass when re-run alone at low load.

## Links

### Out

- verifies [inbox-omits-resumable-landed-run](../capability/rhiza/orchestration/inbox/inbox-omits-resumable-landed-run.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
