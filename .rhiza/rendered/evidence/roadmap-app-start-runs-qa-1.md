# roadmap-app-start-runs-qa-1

evidence · `01M6765RAEKRHYBQ0XJDX16HYH` · actor [claude-code-a83eb6cd](../actor/claude-code-a83eb6cd.md)

## result

PASS

## statement

Built binary in a real Rex tab, with a logging rex shim on PATH that records 'window new' instead of starting a real run and forwards other calls to real Rex. p showed 'Start an interactive plan run ... y/n'; n called nothing; p then y called 'rex window new -s <session> 01M675WMC -c <repo> --json -- <rhiza> run start <cs> --workflow plan' (no --autonomous) and said 'Opened Rex tab 01M675WMC running the plan run'; s then y called the same with 'run start <cs> --autonomous --permission-mode auto'; the app stayed open on the list. o with no run said 'No live run for <cs>'; o with a live (flock-held) run journal whose process names a real Rex block ran real 'rex block inspect' and 'rex focus', and Rex's active window became that tab ('Focused Rex tab block:... of run ...'); after closing that block, o said 'The tab of run ... is gone'.

## Links

### Out

- verifies [roadmap-app-start-runs](../capability/rhiza/roadmap-app-start-runs.md)

## History

- `01M673223RSVC166NJVVA13X75` The roadmap app is a cursor list of milestones and change sets that the owner acts on: open a change set's detail, plan it, ship it, accept it, answer its blocking questions, or jump to its run's Rex tab
