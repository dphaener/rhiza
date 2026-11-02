# plan-requests-qa-2

evidence · `01M6D6G8R4MB30H4EB2HKMWQKT` · actor [claude-code-12fc4e0f](../actor/claude-code-12fc4e0f.md)

## result

PASS

## statement

Built rhiza from the branch head and ran it in a scratch repo with a plan workflow, using a rex stub that only logs window new: request ls says No plan requests; blank text and --kind chore exit 2 with usage errors; a feature, a bug and a milestone queue and show queued; --dry-run writes nothing. Paused, daemon --once starts nothing and opens no tab. Unpaused with --max-runs 1 it opens one tab for the oldest request: run start --workflow plan --new Add CSV exports --request <id> --permission-mode manual, not autonomous, and request ls shows planning with the run; a second --once with one slot starts nothing while that start is unjournaled. With rhiza.autonomy.plan auto and --max-runs 3 it starts the bug as --new Fix: Login fails on Safari ... --autonomous and the milestone with --milestone --autonomous. run start --request refuses an unknown request and needs --new. I could not run the plan run by hand, because run start refuses agent sessions, so the planned state with its change set comes from plan-requests.txtar passing (go test, -count=1) along with describe-golden.

## Links

### Out

- verifies [plan-requests](../capability/rhiza/orchestration/daemon/plan-requests.md)

## History

- `01M6CRZQYG5PJV7X2B234DJCV4` Feature, bug and milestone requests queue in the clone and the daemon turns each into a plan run
