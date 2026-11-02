# request-queue-qa-1

evidence · `01M6DA9083YDDY00CFXWKAEEBN` · actor [claude-code-5c3b0cf5](../actor/claude-code-5c3b0cf5.md)

## result

PASS

## statement

At ca703146 I built the binary outside PATH and used it in a scratch repo with the run's RHIZA_* env unset and a logging-only fake rex (the daemon's tabs are not real agents). rhiza request ls printed 'No plan requests'; an empty text and --kind chore exited 2 with clear errors; a feature, a bug and a milestone were queued into .git/rhiza/requests.json and shown queued. With rhiza.daemon.paused true, daemon --once started nothing and opened no tab. Unpaused, --max-runs 1 --once started only the oldest: rex window new ... run start --workflow plan --new Add CSV export --request <id> --permission-mode manual (no --autonomous), and request ls showed it planning with its run. After daemon set rhiza.autonomy.plan auto, --max-runs 3 --once started the other two with --autonomous: 'Fix: Crash on empty repo' for the bug and --milestone for the milestone. Running the first request's run start by hand (fake claude) succeeded and request ls showed it planned with its run and change set 01M6DA6Q; a second start was refused 'already has run'; starting a still-queued request was refused 'is queued'. In a real Rex window, rhiza tui app roadmap: n then 'Add JSON import' + enter and m then 'Version 3' + enter showed 'Queued a ... request ... for the daemon to plan', request ls listed both queued, and no rex tab was opened. go test plan-requests and describe-golden passed.

## Links

### Out

- verifies [request-queue](../capability/rhiza/orchestration/daemon/plan-requests/request-queue.md)

## History

- `01M6D4SB0D8F9SPGK471DF9RV6` The request queue and the daemon turning requests into plan runs, as a leaf of plan requests
